package chip8

import (
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

func (d *Display) Clear() {
	d.win.Clear(colornames.Black)
	d.win.Update()
}
