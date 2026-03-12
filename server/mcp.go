package server

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"tts-mcp/audio"
	"tts-mcp/tts"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type App struct {
	TTS          tts.TTSProvider
	AudioPlayer  *audio.Player
	DefaultVoice string
	mu           sync.Mutex
}

func NewMCPServer(app *App) *server.MCPServer {
	s := server.NewMCPServer(
		"tts-mcp",
		"0.1.0",
		server.WithToolCapabilities(true),
	)
	s.AddTool(talkTool(app.TTS.ValidVoices()), app.handleTalk)
	return s
}

func talkTool(voices []string) mcp.Tool {
	voiceList := strings.Join(voices, ", ")
	return mcp.NewTool("talk",
		mcp.WithDescription(
			"Speak text aloud using text-to-speech. Converts text to speech and plays it "+
				"through the system's default audio device. Blocks until playback finishes. "+
				"Long text is automatically split into chunks.",
		),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("The text to speak aloud."),
		),
		mcp.WithString("voice",
			mcp.Description(fmt.Sprintf("Voice override. Available: %s. Defaults to server setting.", voiceList)),
		),
	)
}

func (a *App) handleTalk(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, err := req.RequireString("text")
	if err != nil {
		return mcp.NewToolResultError("text is required"), nil
	}
	if strings.TrimSpace(text) == "" {
		return mcp.NewToolResultError("text cannot be empty"), nil
	}

	voice := req.GetString("voice", a.DefaultVoice)
	voice = strings.ToLower(voice)

	valid := false
	for _, v := range a.TTS.ValidVoices() {
		if v == voice {
			valid = true
			break
		}
	}
	if !valid {
		return mcp.NewToolResultError(fmt.Sprintf("unknown voice %q — valid: %s", voice, strings.Join(a.TTS.ValidVoices(), ", "))), nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	chunks := tts.ChunkText(text, a.TTS.MaxTextLength())

	for i, chunk := range chunks {
		log.Printf("Synthesizing chunk %d/%d (%d chars, voice=%s)", i+1, len(chunks), len(chunk), voice)

		body, err := a.TTS.Synthesize(ctx, chunk, voice)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("TTS API error: %v", err)), nil
		}

		if err := a.AudioPlayer.Play(ctx, body); err != nil {
			body.Close()
			return mcp.NewToolResultError(fmt.Sprintf("audio playback error: %v", err)), nil
		}
		body.Close()
	}

	return mcp.NewToolResultText(fmt.Sprintf("Spoke %d characters using voice '%s'.", len(text), voice)), nil
}
