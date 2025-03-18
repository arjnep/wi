package main

import (
	"os"
	"strings"

	"github.com/gdamore/tcell"
)

func (e *Editor) commandMode() {
	command := ""
	e.drawCommandLine(":")

	for {
		ev := e.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEnter {
				go e.executeCommand(command)
				return
			} else if ev.Key() == tcell.KeyEscape {
				return
			} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				if len(command) > 0 {
					command = command[:len(command)-1]
				}
			} else if ev.Rune() != 0 {
				command += string(ev.Rune())
			}
			e.drawCommandLine(":" + command)
		}
	}

}

func (e *Editor) executeCommand(command string) {
	parts := strings.Fields(command) // Split the command by spaces

	switch parts[0] {
	case "w":
		if len(parts) > 1 { // If a filename is provided
			e.filename = parts[1] // Set the new filename
		}
		if e.filename != "" {
			e.saveFile()
		} else {
			e.drawCommandLine("No filename. Use ':w filename'")
		}
	case "q":
		e.screen.Fini()
		os.Exit(0)
	case "wq":
		if e.filename != "" {
			e.saveFile()
		} else {
			e.drawCommandLine("No filename. Use ':w filename'")
			return
		}
		e.screen.Fini()
		os.Exit(0)
	}
}
