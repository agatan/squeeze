package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

const version = "0.0.1"

type mode uint

const (
	normal mode = iota
	regex
	fuzzy
)

var currentMode = normal

func printVersion() {
	fmt.Printf("squeeze - version %s\n", version)
}

func main() {
	// arguments
	var v = flag.Bool("v", false, "show version")
	var re = flag.Bool("re", false, "use regex pattern")
	flag.BoolVar(v, "version", false, "show version")
	flag.BoolVar(re, "regex", false, "use regex pattern")

	flag.Parse()

	if *v {
		printVersion()
		os.Exit(0)
	}
	if *re {
		currentMode = regex
	}

	s, err := newScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	result := ""
	defer func() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		termbox.Close()
		if result != "" {
			fmt.Println(result)
		}
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
				result = s.getSelectedLine().str
				return
			case termbox.KeyCtrlSlash:
				*re = !*re
				updateWithFilter = true
			case termbox.KeySpace:
				s.insertChar(' ')
				updateWithFilter = true
			default:
				s.insertChar(ev.Ch)
				updateWithFilter = true
			}
		case termbox.EventResize:
			s.width, s.height = termbox.Size()
			updateAll = true
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
