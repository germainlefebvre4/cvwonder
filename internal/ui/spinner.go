package ui

import (
	"bytes"
	"time"

	"github.com/charmbracelet/colorprofile"
)

const (
	spinnerDelay    = 200 * time.Millisecond
	spinnerInterval = 80 * time.Millisecond
)

// timers abstracts time so tests can drive the spinner deterministically.
type timers struct {
	// after delivers once after d.
	after func(d time.Duration) <-chan time.Time
	// tick delivers every d until the returned stop function is called.
	tick func(d time.Duration) (<-chan time.Time, func())
}

func realTimers() timers {
	return timers{
		after: time.After,
		tick: func(d time.Duration) (<-chan time.Time, func()) {
			t := time.NewTicker(d)
			return t.C, t.Stop
		},
	}
}

// spinner redraws one line in place while a span runs. It starts only after a
// short delay so fast steps never flicker.
type spinner struct {
	sw      *SharedWriter
	profile colorprofile.Profile
	text    string
	timers  timers
	stop    chan struct{}
	done    chan struct{}
}

func startSpinner(sw *SharedWriter, p colorprofile.Profile, text string, tm timers) *spinner {
	s := &spinner{sw: sw, profile: p, text: text, timers: tm, stop: make(chan struct{}), done: make(chan struct{})}
	// Created before the goroutine starts so a fake clock is armed on return.
	after := tm.after(spinnerDelay)
	go s.run(after)
	return s
}

func (s *spinner) frame(i int) []byte {
	var buf bytes.Buffer
	w := &colorprofile.Writer{Forward: &buf, Profile: s.profile}
	_, _ = w.WriteString(styleAccent.Render(spinnerFrames[i%len(spinnerFrames)]) + " " + s.text)
	return buf.Bytes()
}

func (s *spinner) run(after <-chan time.Time) {
	defer close(s.done)
	select {
	case <-s.stop:
		return
	case <-after:
	}
	tick, stopTick := s.timers.tick(spinnerInterval)
	defer stopTick()
	i := 0
	s.sw.setSpinner(s.frame(i))
	for {
		select {
		case <-s.stop:
			return
		case <-tick:
			i++
			s.sw.setSpinner(s.frame(i))
		}
	}
}

// finish stops the animation and erases its line, restoring the cursor.
func (s *spinner) finish() {
	close(s.stop)
	<-s.done
	s.sw.clearSpinner()
}
