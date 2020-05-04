package main

import (
	"path/filepath"

	"github.com/morinokami/go-chip8/chip8"
)

func main() {
	display := chip8.NewDisplay()
	emulator := chip8.New(display)

	// temp
	path, err := filepath.Abs("./chip8/IBM.ch8")
	if err != nil {
		panic(err)
	}

	emulator.Load(path)
	display.Run(func() {
		display.Init()
		emulator.Run()
	})
}
