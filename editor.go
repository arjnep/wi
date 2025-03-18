package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell"
)

const (
	MODE_VIEW = iota
	MODE_EDIT
)

type Editor struct {
	screen    tcell.Screen
	mode      int
	buffer    [][]rune
	cursorX   int
	cursorY   int
	clipboard []rune
	undoStack []EditorState
	redoStack []EditorState
	filename  string
}

func InitEditor() *Editor {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Error creating screen: %v", err)
	}
	err = screen.Init()
	if err != nil {
		log.Fatalf("Error initializing screen: %v", err)
	}

	var filename string
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	editor := &Editor{
		screen:    screen,
		mode:      MODE_VIEW,
		buffer:    [][]rune{[]rune{}},
		cursorX:   0,
		cursorY:   0,
		undoStack: []EditorState{},
		redoStack: []EditorState{},
	}

	if filename != "" {
		editor.filename = filename
		editor.loadFile(filename)
	}

	return editor
}

func (e *Editor) Run() {
	defer e.screen.Fini()
	if e.filename == "" {
		e.greet()
	}
	e.screen.Clear()

	for {
		e.draw()
		ev := e.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Rune() == 'q' && e.mode == MODE_VIEW {
				return
			}
			e.handleKeyEvent(ev)
		}

	}

}
