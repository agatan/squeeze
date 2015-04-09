package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func main() {
	s := newScreen()

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	result := ""
	defer func() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		termbox.Close()
		fmt.Println(result)
	}()

	s.drawScreen()
	for {
		updatePrompt := false
		updateAll := false
		updateWithFilter := false
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlD:
				return
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				s.moveToLeft()
				updatePrompt = true
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				s.moveToRight()
				updatePrompt = true
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				s.selectNext()
				updateAll = true
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				s.selectPrev()
				updateAll = true
			case termbox.KeyCtrlE, termbox.KeyEnd:
				s.moveToEnd()
				updatePrompt = true
			case termbox.KeyCtrlA, termbox.KeyHome:
				s.moveToBegin()
				updatePrompt = true
			case termbox.KeyDelete, termbox.KeyCtrlH, termbox.KeyBackspace2:
				s.deleteChar()
				updateWithFilter = true
			case termbox.KeyEnter:
				result = s.getSelectedLine()
				return
			default:
				s.insertChar(ev.Ch)
				updateWithFilter = true
			}
		}
		if updateAll {
			s.drawScreen()
		}
		if updatePrompt {
			s.drawPrompt()
		}
		if updateWithFilter {
			updateFilterAndShow(s)
		}
	}
}

func updateFilterAndShow(s *screen) {
	s.drawPrompt()
	go func() {
		result := filtering(s.candidates, s.input)
		s.filtered = <-result
		s.drawScreen()
	}()
}
