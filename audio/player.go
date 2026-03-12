package audio

import (
	"context"
	"io"
	"time"

	"github.com/ebitengine/oto/v3"
)

type Player struct {
	ctx *oto.Context
}

func NewPlayer() (*Player, error) {
	op := &oto.NewContextOptions{
		SampleRate:   24000,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	}
	ctx, readyCh, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}
	<-readyCh
	return &Player{ctx: ctx}, nil
}

// Play streams PCM data from the reader to the default audio device.
// Blocks until playback finishes or the context is cancelled.
func (p *Player) Play(ctx context.Context, r io.Reader) error {
	player := p.ctx.NewPlayer(r)
	player.Play()

	done := make(chan struct{})
	go func() {
		for player.IsPlaying() {
			time.Sleep(10 * time.Millisecond)
		}
		close(done)
	}()

	select {
	case <-done:
		return player.Close()
	case <-ctx.Done():
		player.Close()
		return ctx.Err()
	}
}
