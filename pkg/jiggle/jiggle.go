package jiggle

import (
	"context"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
)

type Jigglr struct {
	mu      sync.Mutex
	running bool
	cancel  context.CancelFunc
	period  time.Duration
}

func New(period time.Duration) *Jigglr {
	if period <= 0 {
		period = 30 * time.Second
	}
	return &Jigglr{period: period}
}

func (j *Jigglr) IsRunning() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.running
}

func (j *Jigglr) Period() time.Duration {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.period
}

func (j *Jigglr) SetPeriod(period time.Duration) {
	if period <= 0 {
		return
	}

	j.mu.Lock()
	wasRunning := j.running
	j.mu.Unlock()

	if wasRunning {
		j.Stop()
		j.Start(period)
		return
	}

	j.mu.Lock()
	j.period = period
	j.mu.Unlock()
}

func (j *Jigglr) Start(period time.Duration) {
	if period <= 0 {
		return
	}

	j.mu.Lock()
	if j.running {
		j.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	j.cancel = cancel
	j.running = true
	j.period = period
	j.mu.Unlock()

	go func() {
		t := time.NewTicker(period)
		defer t.Stop()

		// optional: nudge immediately on start
		nudge()

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				nudge()
			}
		}
	}()
}

func (j *Jigglr) Stop() {
	j.mu.Lock()
	if !j.running {
		j.mu.Unlock()
		return
	}
	cancel := j.cancel
	j.cancel = nil
	j.running = false
	j.mu.Unlock()

	if cancel != nil {
		cancel()
	}
}

func nudge() {
	robotgo.MoveRelative(0, -1)
	robotgo.MoveRelative(0, 1)
}
