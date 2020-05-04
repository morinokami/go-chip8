package chip8

import "time"

type Emulator struct {
	display *Display
	buffer  []byte
}

func New(d *Display) *Emulator {
	return &Emulator{display: d}
}

func (e *Emulator) Run() {
	t := time.NewTicker(16 * time.Millisecond)
	for {
		select {
		case <-t.C:
			e.Cycle()
			e.display.Render(e.buffer)
		}
	}
}

func (e *Emulator) Cycle() {
	// fetch -> decode -> execute
}
