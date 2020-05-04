package main

import "github.com/morinokami/go-chip8/chip8"

func main() {
	display := chip8.NewDisplay()
	emulator := chip8.New(display)
	display.Run(func() {
		display.Init()
		emulator.Run()
	})
}
