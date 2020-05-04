package chip8

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	displayWidth  = 640
	displayHeight = 320
	scale         = 10
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
		Bounds: pixel.R(0, 0, displayWidth, displayHeight),
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

func (d *Display) Render(buf []byte) {
	d.win.Clear(colornames.Black)

	imd := imdraw.New(nil)
	x := 1.0
	y := 2.0
	imd.Color = colornames.Pink
	imd.Push(pixel.V(x*scale, displayHeight-(y*scale)), pixel.V((x+1)*scale, displayHeight-((y+1)*scale)))
	imd.Rectangle(0)
	x = 63.0
	y = 31.0
	imd.Color = colornames.Pink
	imd.Push(pixel.V(x*scale, displayHeight-(y*scale)), pixel.V((x+1)*scale, displayHeight-((y+1)*scale)))
	imd.Rectangle(0)
	imd.Draw(d.win)

	d.win.Update()
}
