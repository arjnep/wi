package main

type EditorState struct {
	buffer  [][]rune
	cursorX int
	cursorY int
}

func (e *Editor) copyLine() {
	if len(e.buffer) > 0 {
		e.clipboard = append([]rune{}, e.buffer[e.cursorY]...)
	}
}

func (e *Editor) cutLine() {

	if len(e.buffer) > 0 {
		e.clipboard = append([]rune{}, e.buffer[e.cursorY]...)

		if len(e.buffer) == 1 {
			e.buffer = [][]rune{[]rune{}}
			e.cursorY = 0
			e.cursorX = 0
		} else {
			e.buffer = append(e.buffer[:e.cursorY], e.buffer[e.cursorY+1:]...)

			if e.cursorY >= len(e.buffer) {
				e.cursorY = max(0, len(e.buffer)-1)
			}

			if e.cursorX > len(e.buffer[e.cursorY]) {
				e.cursorX = len(e.buffer[e.cursorY])
			}
		}
	}
}

func (e *Editor) pasteLine() {

	if len(e.clipboard) > 0 {
		newLine := append([]rune{}, e.clipboard...)
		e.buffer = append(e.buffer[:e.cursorY+1], append([][]rune{newLine}, e.buffer[e.cursorY+1:]...)...)
		e.cursorY++
		e.cursorX = len(e.clipboard) // To change or not to change thats the question
	}
}

func (e *Editor) saveState() {
	stateCopy := make([][]rune, len(e.buffer))
	for i := range e.buffer {
		stateCopy[i] = append([]rune{}, e.buffer[i]...)
	}

	// push current state to undo stack
	e.undoStack = append(e.undoStack, EditorState{
		buffer:  stateCopy,
		cursorX: e.cursorX,
		cursorY: e.cursorY,
	})
	e.redoStack = nil // clear redo stack on new change
}

func (e *Editor) undo() {
	if len(e.undoStack) == 0 {
		return
	}

	e.redoStack = append(e.redoStack, EditorState{
		buffer:  e.buffer,
		cursorX: e.cursorX,
		cursorY: e.cursorY,
	})

	lastState := e.undoStack[len(e.undoStack)-1]
	e.undoStack = e.undoStack[:len(e.undoStack)-1]

	e.buffer = lastState.buffer
	e.cursorX = min(lastState.cursorX, len(e.buffer[lastState.cursorY]))
	e.cursorY = min(lastState.cursorY, len(e.buffer)-1)
}

func (e *Editor) redo() {
	if len(e.redoStack) == 0 {
		return
	}

	e.undoStack = append(e.undoStack, EditorState{
		buffer:  e.buffer,
		cursorX: e.cursorX,
		cursorY: e.cursorY,
	})

	lastState := e.redoStack[len(e.redoStack)-1]
	e.redoStack = e.redoStack[:len(e.redoStack)-1]

	e.buffer = lastState.buffer
	e.cursorX = min(lastState.cursorX, len(e.buffer[lastState.cursorY]))
	e.cursorY = min(lastState.cursorY, len(e.buffer)-1)
}
