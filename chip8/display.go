package chip8

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Display struct {
	win *pixelgl.Window
}

func NewDisplay() *Display {
	return &Display{}
}

func (d *Display) Init() {
	cfg := pixelgl.WindowConfig{
		Title:  "CHIP-8",
		Bounds: pixel.R(0, 0, DisplayWidth, DisplayHeight),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	d.win = win
}

func (d *Display) Run(f func()) {
	pixelgl.Run(f)
}

func (d *Display) Render(buf [BufferSize]byte) {
	d.win.Clear(colornames.Black)

	imd := imdraw.New(nil)
	for i, b := range buf {
		if b == 1 {
			x := float64(i % 64)
			y := float64(i / 64)
			imd.Color = colornames.Pink
			imd.Push(
				pixel.V(x*ScalingFactor, DisplayHeight-(y*ScalingFactor)),
				pixel.V((x+1)*ScalingFactor, DisplayHeight-((y+1)*ScalingFactor)),
			)
			imd.Rectangle(0)
		}
	}
	imd.Draw(d.win)

	d.win.Update()
}

func (d *Display) Key() (byte, bool) {
	if d.win.Pressed(pixelgl.Key4) {
		return 0x1, true
	} else if d.win.Pressed(pixelgl.Key5) {
		return 0x2, true
	} else if d.win.Pressed(pixelgl.Key6) {
		return 0x3, true
	} else if d.win.Pressed(pixelgl.Key7) {
		return 0xC, true
	} else if d.win.Pressed(pixelgl.KeyR) {
		return 0x4, true
	} else if d.win.Pressed(pixelgl.KeyT) {
		return 0x5, true
	} else if d.win.Pressed(pixelgl.KeyY) {
		return 0x6, true
	} else if d.win.Pressed(pixelgl.KeyU) {
		return 0xD, true
	} else if d.win.Pressed(pixelgl.KeyF) {
		return 0x7, true
	} else if d.win.Pressed(pixelgl.KeyG) {
		return 0x8, true
	} else if d.win.Pressed(pixelgl.KeyH) {
		return 0x9, true
	} else if d.win.Pressed(pixelgl.KeyJ) {
		return 0xE, true
	} else if d.win.Pressed(pixelgl.KeyV) {
		return 0xA, true
	} else if d.win.Pressed(pixelgl.KeyB) {
		return 0x0, true
	} else if d.win.Pressed(pixelgl.KeyN) {
		return 0xB, true
	} else if d.win.Pressed(pixelgl.KeyM) {
		return 0xF, true
	} else {
		return 0x0, false
	}
}

func (d *Display) Beep() {
	fmt.Println("Beep!")
}
