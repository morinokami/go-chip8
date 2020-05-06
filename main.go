package main

import (
	"errors"
	"log"
	"os"

	"github.com/morinokami/go-chip8/chip8"
	"github.com/morinokami/go-chip8/games"
	"github.com/urfave/cli/v2"
)

func main() {
	display := chip8.NewDisplay()
	emulator := chip8.New(display)

	var game int
	app := &cli.App{
		Name:  "go-chip8",
		Usage: "a CHIP-8 emulator written in Go",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "game",
				Aliases:     []string{"g"},
				Usage:       "enter one of the following numbers: " + games.AvailableGames(),
				Destination: &game,
			},
		},
		Action: func(c *cli.Context) error {
			if game < 0 || len(games.Games)-1 < game {
				return errors.New("invalid game id")
			}

			emulator.Load(games.Games[game].Binary)
			display.Run(func() {
				display.Init()
				emulator.Run()
			})

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
