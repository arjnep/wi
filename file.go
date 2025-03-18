package main

import (
	"bufio"
	"os"
)

func (e *Editor) loadFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	e.buffer = [][]rune{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		e.buffer = append(e.buffer, []rune(scanner.Text()))
	}

	if len(e.buffer) == 0 {
		e.buffer = [][]rune{[]rune{}}
	}
}

func (e *Editor) saveFile() {
	if e.filename == "" {
		e.drawCommandLine("No filename. Use ':w <filename>'")
		return
	}

	file, err := os.Create(e.filename)
	if err != nil {
		e.drawCommandLine("Error saving file")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range e.buffer {
		writer.WriteString(string(line) + "\n")
	}
	writer.Flush()
	e.drawCommandLine("Saved: " + e.filename)
}
