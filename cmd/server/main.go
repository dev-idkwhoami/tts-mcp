package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"tts-mcp/audio"
	mcpserver "tts-mcp/server"
	"tts-mcp/tts"

	goserver "github.com/mark3labs/mcp-go/server"
)

func main() {
	log.SetFlags(log.Ltime)
	log.SetOutput(os.Stderr)

	provider := flag.String("provider", "", "TTS provider: xai, openai (env: PROVIDER, default: xai)")
	key := flag.String("key", "", "API key (env: API_KEY)")
	voice := flag.String("voice", "", "Default TTS voice (provider-specific, see README)")
	lang := flag.String("lang", "en", "Language code for xAI (e.g. en, de, fr, ja)")
	flag.Parse()

	// Resolve provider: flag → env → default
	p := strings.ToLower(*provider)
	if p == "" {
		p = strings.ToLower(os.Getenv("PROVIDER"))
	}
	if p == "" {
		p = "xai"
	}

	// Resolve API key: flag → env → fatal
	apiKey := *key
	if apiKey == "" {
		apiKey = os.Getenv("API_KEY")
	}
	if apiKey == "" {
		log.Fatal("API key required: set -key flag or API_KEY env var")
	}

	// Build TTS provider
	var ttsP tts.TTSProvider
	switch p {
	case "xai":
		ttsP = tts.NewXAI(apiKey, *lang)
	case "openai":
		ttsP = tts.NewOpenAI(apiKey)
	default:
		log.Fatalf("Unknown provider %q — valid: xai, openai", p)
	}

	// Resolve default voice
	v := strings.ToLower(*voice)
	if v == "" {
		v = ttsP.ValidVoices()[0]
	} else {
		valid := false
		for _, vv := range ttsP.ValidVoices() {
			if vv == v {
				valid = true
				break
			}
		}
		if !valid {
			log.Fatalf("Invalid voice %q for %s — valid: %s", *voice, p, strings.Join(ttsP.ValidVoices(), ", "))
		}
	}

	// Init audio
	player, err := audio.NewPlayer()
	if err != nil {
		log.Fatalf("Failed to initialize audio player: %v", err)
	}

	app := &mcpserver.App{
		TTS:          ttsP,
		AudioPlayer:  player,
		DefaultVoice: v,
	}

	log.Printf("Starting tts-mcp (provider=%s, voice=%s)", p, v)
	mcpSrv := mcpserver.NewMCPServer(app)
	stdio := goserver.NewStdioServer(mcpSrv)

	ctx := context.Background()
	if err := stdio.Listen(ctx, os.Stdin, os.Stdout); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}
