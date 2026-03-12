package tts

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type XAI struct {
	apiKey     string
	language   string
	httpClient *http.Client
}

type xaiRequest struct {
	Text         string           `json:"text"`
	VoiceID      string           `json:"voice_id"`
	Language     string           `json:"language"`
	OutputFormat *xaiOutputFormat `json:"output_format,omitempty"`
}

type xaiOutputFormat struct {
	Codec      string `json:"codec"`
	SampleRate int    `json:"sample_rate"`
}

func NewXAI(apiKey, language string) *XAI {
	return &XAI{
		apiKey:   apiKey,
		language: language,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (x *XAI) MaxTextLength() int { return 15000 }

func (x *XAI) ValidVoices() []string {
	return []string{"eve", "ara", "rex", "sal", "leo"}
}

func (x *XAI) Synthesize(ctx context.Context, text, voice string) (io.ReadCloser, error) {
	reqBody := xaiRequest{
		Text:     text,
		VoiceID:  voice,
		Language: x.language,
		OutputFormat: &xaiOutputFormat{
			Codec:      "pcm",
			SampleRate: 24000,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	return doTTSRequest(ctx, x.httpClient, "https://api.x.ai/v1/tts", x.apiKey, bodyBytes)
}
