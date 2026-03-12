package tts

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// doTTSRequest performs an HTTP POST with retry logic shared by all TTS providers.
func doTTSRequest(ctx context.Context, client *http.Client, url, apiKey string, body []byte) (io.ReadCloser, error) {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			delay := time.Duration(1<<uint(attempt)) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		if resp.StatusCode == 200 {
			return resp.Body, nil
		}

		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		resp.Body.Close()

		switch resp.StatusCode {
		case 400:
			return nil, fmt.Errorf("bad request: %s", string(errBody))
		case 401:
			return nil, fmt.Errorf("unauthorized — check your API key")
		case 429, 500, 503:
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(errBody))
			continue
		default:
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(errBody))
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}
