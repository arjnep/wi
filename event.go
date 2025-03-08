package main

import (
	"github.com/gdamore/tcell"
)

func (e *Editor) insertCharacter(ev *tcell.EventKey) {

	if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
		if e.cursorX > 0 {
			e.buffer[e.cursorY] = append(e.buffer[e.cursorY][:e.cursorX-1], e.buffer[e.cursorY][e.cursorX:]...)
			e.cursorX--
		} else if e.cursorY > 0 {
			prevLen := len(e.buffer[e.cursorY-1])
			e.buffer[e.cursorY-1] = append(e.buffer[e.cursorY-1], e.buffer[e.cursorY]...)
			e.buffer = append(e.buffer[:e.cursorY], e.buffer[e.cursorY+1:]...)
			e.cursorY--
			e.cursorX = prevLen
		}
	} else if ev.Key() == tcell.KeyEnter {
		newLine := e.buffer[e.cursorY][e.cursorX:]
		e.buffer[e.cursorY] = e.buffer[e.cursorY][:e.cursorX]
		e.buffer = append(e.buffer[:e.cursorY+1], append([][]rune{newLine}, e.buffer[e.cursorY+1:]...)...)
		e.cursorY++
		e.cursorX = 0
	} else if ev.Rune() != 0 {
		line := e.buffer[e.cursorY]
		line = append(line[:e.cursorX], append([]rune{ev.Rune()}, line[e.cursorX:]...)...)
		e.buffer[e.cursorY] = line
		e.cursorX++
	}
}

func (e *Editor) handleKeyEvent(ev *tcell.EventKey) {
	if e.mode == MODE_VIEW {
		switch ev.Rune() {
		case 'i':
			e.mode = MODE_EDIT
		case 'y':
			e.copyLine()
		case 'd':
			e.cutLine()
		case 'p':
			e.pasteLine()
		case 'u':
			e.undo()
		case 'r':
			e.redo()
		}
	} else if e.mode == MODE_EDIT {
		if ev.Key() == tcell.KeyEscape {
			e.mode = MODE_VIEW
			return
		}
		e.saveState()
		e.insertCharacter(ev)
	}
	e.handleNavigation(ev)

}

func (e *Editor) handleNavigation(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyLeft:
		e.keyLeft()
	case tcell.KeyRight:
		e.keyRight()
	case tcell.KeyUp:
		e.keyUp()
	case tcell.KeyDown:
		e.keyDown()
	}
}

func (e *Editor) keyLeft() {
	if e.cursorX > 0 {
		e.cursorX--
	} else if e.cursorY > 0 {
		e.cursorY--
		e.cursorX = len(e.buffer[e.cursorY])
	}
}
func (e *Editor) keyRight() {
	if e.cursorX < len(e.buffer[e.cursorY]) {
		e.cursorX++
	} else if e.cursorY < len(e.buffer)-1 {
		e.cursorY++
		e.cursorX = 0
	}
}
func (e *Editor) keyUp() {
	if e.cursorY > 0 {
		e.cursorY--
		if e.cursorX > len(e.buffer[e.cursorY]) {
			e.cursorX = len(e.buffer[e.cursorY])
		}
	}
}
func (e *Editor) keyDown() {
	if e.cursorY < len(e.buffer)-1 {
		e.cursorY++
		if e.cursorX > len(e.buffer[e.cursorY]) {
			e.cursorX = len(e.buffer[e.cursorY])
		}
	}
}
