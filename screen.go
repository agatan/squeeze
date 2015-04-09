package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/andrew-d/go-termutil"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type screen struct {
	width, height int
	cursorX       int
	selectedLine  int
	candidates    []string
	filtered      []string
	input         []rune
}

func newScreen() *screen {
	if termutil.Isatty(os.Stdin.Fd()) {
		fmt.Fprintf(os.Stderr, "nothing to read\n")
		os.Exit(1)
	}
	scanner := bufio.NewScanner(os.Stdin)
	var candidates = []string{}
	for scanner.Scan() {
		candidates = append(candidates, scanner.Text())
	}

	s := new(screen)
	s.width, s.height = termbox.Size()
	s.candidates = candidates
	s.filtered = candidates
	return s
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

func (s *screen) selectNext() {
	if s.selectedLine < len(s.filtered) {
		s.selectedLine++
	}
}

func (s *screen) selectPrev() {
	if s.selectedLine > 0 {
		s.selectedLine--
	}
}

func (s *screen) getSelectedLine() string {
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
	s.setPrompt()
	termbox.Flush()
}

func (s *screen) drawScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	s.setPrompt()
	for idx, str := range s.filtered {
		if idx == s.selectedLine {
			setLine(0, idx+1, termbox.ColorDefault, termbox.ColorGreen, str)
		} else {
			setLine(0, idx+1, termbox.ColorDefault, termbox.ColorDefault, str)
		}
	}
	termbox.Flush()
}
