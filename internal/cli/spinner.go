package cli

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type spinner struct {
	w     io.Writer
	done  chan struct{}
	chars []rune
	wg    sync.WaitGroup
	index int
	once  sync.Once
}

func newSpinner(w io.Writer) *spinner {
	return &spinner{
		w:     w,
		chars: []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'},
		done:  make(chan struct{}),
	}
}

func (s *spinner) start(banner string) {
	s.wg.Go(func() {
		if banner != "" {
			fmt.Fprintln(s.w, banner)
		}

		for {
			select {
			case <-s.done:
				return
			default:
				fmt.Fprintf(s.w, "    %c\r", s.chars[s.index])
				s.index = (s.index + 1) % len(s.chars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	})
}

func (s *spinner) stop() {
	s.once.Do(func() {
		close(s.done)
		// Wait for the last frame, so that clearing the line comes after it.
		s.wg.Wait()

		// Clear spinner line.
		fmt.Fprint(s.w, "\r\033[K")
	})
}
