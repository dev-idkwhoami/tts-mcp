package tts

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAI struct {
	apiKey     string
	httpClient *http.Client
}

type openaiRequest struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format"`
	Speed          float64 `json:"speed"`
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (o *OpenAI) MaxTextLength() int { return 4096 }

func (o *OpenAI) ValidVoices() []string {
	return []string{"alloy", "ash", "coral", "echo", "fable", "nova", "onyx", "sage", "shimmer"}
}

func (o *OpenAI) Synthesize(ctx context.Context, text, voice string) (io.ReadCloser, error) {
	reqBody := openaiRequest{
		Model:          "tts-1",
		Input:          text,
		Voice:          voice,
		ResponseFormat: "pcm",
		Speed:          1.0,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	return doTTSRequest(ctx, o.httpClient, "https://api.openai.com/v1/audio/speech", o.apiKey, bodyBytes)
}
