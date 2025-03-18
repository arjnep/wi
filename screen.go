package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

func (e *Editor) greet() {
	e.screen.Clear()
	greeting := []string{
		"Wi - Not Vi Definitely!",
		"-------------------------------",
		"i = edit mode.",
		"esc = view mode",
		"y = copy line",
		"d = cut line",
		"p = paste line",
		"u = undo",
		"r = redo",
		"q = exit",
		"Press any key to start...",
	}

	w, h := e.screen.Size()
	startY := (h / 2) - (len(greeting) / 2) // centering vertical

	for i, line := range greeting {
		startX := (w / 2) - (len(line) / 2) // centering horizontal
		for j, ch := range line {
			e.screen.SetContent(startX+j, startY+i, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
		}
	}
	e.screen.Show()

	for {
		ev := e.screen.PollEvent()
		if _, ok := ev.(*tcell.EventKey); ok {
			break
		}
	}
	e.screen.Clear()
}

func (e *Editor) draw() {
	e.screen.Clear()
	w, h := e.screen.Size()

	for y := 0; y < h-1 && y < len(e.buffer); y++ {
		e.drawLine(y)
	}

	e.drawStatusBar(w)
	e.screen.Show()
}

func (e *Editor) drawLine(y int) {
	line := e.buffer[y]

	lineNumberWidth := len(fmt.Sprintf("%d", len(e.buffer))) + 1
	lineNumber := fmt.Sprintf("%*d", lineNumberWidth, y+1) // right align
	for x := 0; x < lineNumberWidth; x++ {
		e.screen.SetContent(x, y+1, rune(lineNumber[x]), nil, tcell.StyleDefault.Foreground(tcell.ColorDimGray))
	}

	width, _ := e.screen.Size()
	for x := 0; x < width-lineNumberWidth; x++ {
		var ch rune = ' '

		if x < len(line) {
			ch = line[x]
		}

		lineStyle := tcell.StyleDefault
		if y == e.cursorY {
			lineStyle = tcell.StyleDefault.Background(tcell.ColorLightSlateGray).Foreground(tcell.ColorWhite) // if cursor is on line set the bg gray
		}

		e.screen.SetContent(x+lineNumberWidth+1, y+1, ch, nil, lineStyle)
	}

	if y == e.cursorY {
		e.screen.SetContent(e.cursorX+lineNumberWidth+1, e.cursorY+1, ' ', nil, tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))
	}
}

func (e *Editor) drawCommandLine(message string) {
	width, height := e.screen.Size()

	for x := 0; x < width; x++ {
		e.screen.SetContent(x, height-1, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue))
	}

	for x, ch := range message {
		e.screen.SetContent(x, height-1, ch, nil, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))
	}
	e.screen.Show()
}

func (e *Editor) drawStatusBar(width int) {

	leftStatus := fmt.Sprintf("MODE: %s", e.getMode())
	rightStatus := fmt.Sprintf("Cursor: (%d, %d) | Buffer: [%v][%v]", e.cursorX, e.cursorY, len(e.buffer), len(e.buffer[e.cursorY]))

	rightStartPos := width - len(rightStatus)

	for x := 0; x < width; x++ {
		e.screen.SetContent(x, 0, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue))
	}

	for x, ch := range leftStatus {
		e.screen.SetContent(x+3, 0, ch, nil, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))
	}

	for x, ch := range rightStatus {
		e.screen.SetContent(rightStartPos+x-3, 0, ch, nil, tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite))
	}
}

func (e *Editor) getMode() string {
	if e.mode == MODE_VIEW {
		return "VIEW"
	}
	return "EDIT"
}
