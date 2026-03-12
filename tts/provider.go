package tts

import (
	"context"
	"io"
	"strings"
)

// TTSProvider abstracts text-to-speech synthesis across providers.
type TTSProvider interface {
	Synthesize(ctx context.Context, text, voice string) (io.ReadCloser, error)
	MaxTextLength() int
	ValidVoices() []string
}

// ChunkText splits text into chunks that fit within the given character limit.
// It tries to split at sentence boundaries for natural speech.
func ChunkText(text string, maxLen int) []string {
	if len(text) <= maxLen {
		return []string{text}
	}

	var chunks []string
	for len(text) > 0 {
		if len(text) <= maxLen {
			chunks = append(chunks, text)
			break
		}
		cutPoint := findSentenceBoundary(text, maxLen)
		chunks = append(chunks, text[:cutPoint])
		text = strings.TrimLeft(text[cutPoint:], " \n\r\t")
	}
	return chunks
}

func findSentenceBoundary(text string, maxLen int) int {
	lookback := 1000
	if lookback > maxLen {
		lookback = maxLen
	}
	for i := maxLen; i > maxLen-lookback && i > 0; i-- {
		if text[i-1] == '.' || text[i-1] == '!' || text[i-1] == '?' {
			if i >= len(text) || text[i] == ' ' || text[i] == '\n' {
				return i
			}
		}
	}
	if idx := strings.LastIndex(text[:maxLen], " "); idx > 0 {
		return idx
	}
	return maxLen
}
