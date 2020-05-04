package chip8

import (
	"fmt"
	"io/ioutil"
	"time"
)

type Emulator struct {
	display     *Display
	frameBuffer [BufferSize]byte
	memory      [MemorySize]byte
	pc          int
}

func New(d *Display) *Emulator {
	return &Emulator{display: d, pc: PCStart}
}

func (e *Emulator) Load(rom string) {
	// clear memory
	for i := range e.memory {
		e.memory[i] = 0
	}

	// load fonts into memory
	for i, b := range FontSet {
		e.memory[i] = b
	}

	// load rom into memory
	data, err := ioutil.ReadFile(rom)
	if err != nil {
		panic(err)
	}
	for i, b := range data {
		e.memory[PCStart+i] = b
	}
}

func (e *Emulator) Run() {
	t0 := time.NewTicker(ClockSpeed)
	t1 := time.NewTicker(FrameRate)
	for {
		select {
		case <-t0.C:
			e.Cycle()
		case <-t1.C:
			e.display.Render(e.frameBuffer)
		}
	}
}

func (e *Emulator) Cycle() {
	// fetch -> decode -> execute
	opcode := uint16(e.memory[e.pc])<<8 | uint16(e.memory[e.pc+1])
	e.execute(opcode)
}

func (e *Emulator) decode(opcode uint16) Instruction {
	//kk := opcode & 0x00FF
	//n := opcode & 0x000F

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			return CLS
		case 0x00EE:
			return RET
		default:
			return SYS
		}
	default:
		return UNKNOWN
	}
}

func (e *Emulator) execute(opcode uint16) {
	inst := e.decode(opcode)
	incPC := true

	switch inst {
	case CLS:
		e.display.Clear()
	case UNKNOWN:
		panic(fmt.Sprintf("Unknown opcode: 0x%04x", opcode))
	}

	if incPC {
		e.pc += 2
	}
}
