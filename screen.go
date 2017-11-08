package main

import (
	"bufio"
	"io"
	"os"
	"sync"

	isatty "github.com/mattn/go-isatty"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/pkg/errors"
)

type screen struct {
	rederLock     sync.Mutex
	width, height int
	cursorX       int
	selectedLine  int

	lock       sync.Mutex
	candidates []string
	filtered   []match
	input      []rune
}

func (s *screen) initializeCandidateAsync(from io.Reader) {
	sc := bufio.NewScanner(from)
	ch := make(chan string)
	done := make(chan struct{})
	go func() {
		for sc.Scan() {
			ch <- sc.Text()
		}
		done <- struct{}{}
	}()
	s.appendFromChan(ch, done)
}

func (s *screen) appendFromChan(ch <-chan string, done <-chan struct{}) {
	for {
		select {
		case str := <-ch:
			s.lock.Lock()
			s.candidates = append(s.candidates, str)
			m, err := matching(str, s.input)
			if err == nil {
				s.filtered = append(s.filtered, m)
				updateFilterAndShow(s)
			}
			s.lock.Unlock()
		case <-done:
			return
		}
	}
}

func newScreen() (*screen, error) {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return nil, errors.New("failed to read source")
	}

	if err := termbox.Init(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize terminal state")
	}

	s := new(screen)
	s.width, s.height = termbox.Size()
	go s.initializeCandidateAsync(os.Stdin)
	return s, nil
}

func (s *screen) moveToLeft() {
	if s.cursorX > 0 {
		s.cursorX--
	}
}

func (s *screen) moveToRight() {
	if s.cursorX < len(s.input) {
		s.cursorX++
	}
}

func (s *screen) moveToBegin() {
	s.cursorX = 0
}

func (s *screen) moveToEnd() {
	s.cursorX = len(s.input)
}

func (s *screen) insertChar(ch rune) {
	tmp := []rune{}
	tmp = append(tmp, s.input[0:s.cursorX]...)
	tmp = append(tmp, ch)
	tmp = append(tmp, s.input[s.cursorX:]...)
	s.input = tmp
	s.cursorX++
}

func (s *screen) deleteChar() {
	if s.cursorX == 0 {
		return
	}
	s.input = append(s.input[0:s.cursorX-1], s.input[s.cursorX:]...)
	s.cursorX--
}

func setLine(x, y int, fg, bg termbox.Attribute, strs ...string) {
	for _, str := range strs {
		for _, c := range str {
			termbox.SetCell(x, y, c, fg, bg)
			x += runewidth.RuneWidth(c)
		}
	}
}

func setMatch(y int, fg, bg termbox.Attribute, m match) {
	if len(m.positions) == 0 {
		setLine(0, y, fg, bg, m.str)
		return
	}
	last := 0
	for _, hl := range m.positions {
		setLine(last, y, fg, bg, m.str[last:hl.Start])
		setLine(hl.Start, y, termbox.ColorRed, bg, m.str[hl.Start:hl.End])
		last = hl.End
	}
	setLine(last, y, fg, bg, m.str[last:])
}

func (s *screen) selectNext() {
	if s.selectedLine < len(s.filtered)-1 {
		s.selectedLine++
	}
}

func (s *screen) selectPrev() {
	if s.selectedLine > 0 {
		s.selectedLine--
	}
}

func (s *screen) getSelectedLine() match {
	if len(s.filtered) == 0 {
		return match{}
	}
	return s.filtered[s.selectedLine]
}

func (s *screen) setPrompt() {
	// erase header.
	for x := 0; x < s.width; x++ {
		termbox.SetCell(x, 0, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	prompt := "> "
	setLine(0, 0, termbox.ColorDefault, termbox.ColorDefault, prompt, string(s.input))
	termbox.SetCursor(runewidth.StringWidth(prompt+string(s.input[0:s.cursorX])), 0)
}

func (s *screen) drawPrompt() {
	s.rederLock.Lock()
	defer s.rederLock.Unlock()
	s.setPrompt()
	termbox.Flush()
}

func (s *screen) drawScreen() {
	s.rederLock.Lock()
	defer s.rederLock.Unlock()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	s.setPrompt()
	for idx, m := range s.filtered {
		if idx == s.selectedLine {
			setMatch(idx+1, termbox.ColorDefault, termbox.ColorGreen, m)
		} else {
			setMatch(idx+1, termbox.ColorDefault, termbox.ColorDefault, m)
		}
	}
	termbox.Flush()
}

func updateFilterAndShow(s *screen) {
	s.drawPrompt()
	go func() {
		s.lock.Lock()
		defer s.lock.Unlock()
		var result <-chan []match
		if currentMode == regex {
			result = regexpFiltering(s.candidates, s.input)
		} else {
			result = filtering(s.candidates, s.input)
		}
		s.filtered = <-result
		if s.selectedLine >= len(s.filtered) {
			s.selectedLine = len(s.filtered) - 1
		}
		s.drawScreen()
	}()
}
