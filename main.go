package main

import (
	"fmt"

	"github.com/morinokami/go-chip8/chip8"
	"github.com/morinokami/go-chip8/games"
)

func main() {
	display := chip8.NewDisplay()
	emulator := chip8.New(display)

	fmt.Println(games.AvailableGames())

	emulator.Load(games.Games[8].Binary)
	display.Run(func() {
		display.Init()
		emulator.Run()
	})
}
