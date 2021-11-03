// Copyright (c) 2021, Geert JM Vanderkelen

package xterm

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/geertjanvdk/xkit/xansi"
)

var defaultDelay = time.Millisecond * 500

type spinnerSet struct {
	characters           []string
	notSpinningCharacter string
}

var spinnerSets = map[string]spinnerSet{
	"braille-dots": {
		characters:           []string{"⣾", "⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽"},
		notSpinningCharacter: "\u2800", // BRAILLE PATTERN BLANK
	},
}

type Spinner struct {
	spinnerSet    spinnerSet
	out           io.Writer
	chanStop      chan struct{}
	MessageBefore string        // shown before the spinner
	MessageAfter  string        // shown after the spinner
	Delay         time.Duration // duration between each spin
	Render        *xansi.Render

	spinning bool
	mu       sync.RWMutex
}

func NewSpinner() (*Spinner, error) {
	s := &Spinner{
		spinnerSet: spinnerSets["braille-dots"],
		out:        os.Stdout,
		chanStop:   make(chan struct{}),
		Delay:      defaultDelay,
	}
	return s, nil
}

func (s *Spinner) write(spinChar int, before, after string) {
	var c string
	if spinChar == -1 {
		c = s.spinnerSet.notSpinningCharacter
	} else {
		c = s.spinnerSet.characters[spinChar]
	}

	if s.Render != nil {
		_, _ = fmt.Fprintf(s.out, "\r%s%s%s%s", before,
			s.Render.Sprintf("%s", c), xansi.Reset(), after)
	} else {
		_, _ = fmt.Fprintf(s.out, "\r%s%s%s", before, c, after)
	}
}

func (s *Spinner) Start() {
	s.mu.Lock()
	if s.spinning {
		s.mu.Unlock()
		return
	}
	s.spinning = true
	_, _ = fmt.Fprintf(s.out, xansi.HideCursor)
	s.mu.Unlock()

	go func() {
		for {
			for i := 0; i < len(s.spinnerSet.characters); i++ {
				select {
				case <-s.chanStop:
					return
				default:
					s.mu.Lock()
					s.clearLine()
					s.write(i, s.MessageBefore, s.MessageAfter)
					s.mu.Unlock()
					time.Sleep(s.Delay)
				}
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.StopCustom("", "")
}

func (s *Spinner) StopCustom(before, after string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.spinning {
		s.spinning = false
		s.clearLine()
		if after != "" {
			if before == "" {
				before = s.MessageBefore
			}
			s.write(-1, before, after)
		}
		s.chanStop <- struct{}{}
	}
}

func (s *Spinner) clearLine() {
	_, _ = fmt.Fprintf(s.out, xansi.ClearLineAtCursor)
}
