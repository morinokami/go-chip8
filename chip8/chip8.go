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
	pc          uint16
	iReg        uint16
	vReg        [VRegisterSize]byte
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
	case 0x1000:
		return JPAddr
	case 0x6000:
		return LDVxByte
	case 0x7000:
		return ADDVxByte
	case 0xA000:
		return LDIAddr
	case 0xD000:
		return DRW
	default:
		return UNKNOWN
	}
}

func (e *Emulator) execute(opcode uint16) {
	inst := e.decode(opcode)
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	nnn := opcode & 0x0FFF
	kk := opcode & 0x00FF
	n := opcode & 0x000F
	incPC := true

	switch inst {
	case CLS:
		// 00E0 - CLS
		// Clear the display.
		for i := 0; i < BufferSize; i++ {
			e.frameBuffer[i] = 0
		}
	case JPAddr:
		// 1nnn - JP addr
		// Jump to location nnn.
		//
		// The interpreter sets the program counter to nnn.
		e.pc = nnn
		incPC = false
	case LDVxByte:
		// 6xkk - LD Vx, byte
		// Set Vx = kk.
		//
		// The interpreter puts the value kk into register Vx.
		e.vReg[x] = byte(kk)
	case ADDVxByte:
		// 6xkk - LD Vx, byte
		// Set Vx = kk.
		//
		// The interpreter puts the value kk into register Vx.
		e.vReg[x] = byte((uint16(e.vReg[x]) + kk) & 0xFF)
	case LDIAddr:
		// Annn - LD I, addr
		// Set I = nnn.
		//
		// The value of register I is set to nnn.
		e.iReg = nnn
	case DRW:
		// Dxyn - DRW Vx, Vy, nibble
		// Display n-byte sprite starting at memory location I at (Vx, Vy), set
		// VF = collision.
		//
		// The interpreter reads n bytes from memory, starting at the address
		// stored in I. These bytes are then displayed as sprites on screen at
		// coordinates (Vx, Vy). Sprites are XORed onto the existing screen. If
		// this causes any pixels to be erased, VF is set to 1, otherwise it is
		// set to 0. If the sprite is positioned so part of it is outside the
		// coordinates of the display, it wraps around to the opposite side of
		// the screen. See instruction 8xy3 for more information on XOR, and
		// section 2.4, Display, for more information on the Chip-8 screen and
		// sprites.
		vx := e.vReg[x]
		vy := e.vReg[y]
		sprite := e.memory[e.iReg : e.iReg+n]
		erased := e.drawSprite(vx, vy, sprite)
		if erased {
			e.vReg[0xF] = 1
		} else {
			e.vReg[0xF] = 0
		}
	case UNKNOWN:
		panic(fmt.Errorf("unknown opcode: 0x%04x", opcode))
	}

	if incPC {
		e.pc += 2
	}
}

func (e *Emulator) drawSprite(vx, vy byte, sprite []byte) bool {
	erased := false
	for y, b := range sprite {
		for x, bit := range bits(b) {
			erased = e.drawPixel(x+int(vx), y+int(vy), bit == 1)
		}
	}
	return erased
}

func (e *Emulator) drawPixel(x, y int, fill bool) bool {
	prevFilled := e.filled(x, y)
	curFilled := fill != prevFilled
	if curFilled {
		e.frameBuffer[x+y*BaseWidth] = 1
	} else {
		e.frameBuffer[x+y*BaseWidth] = 0
	}
	return prevFilled && !curFilled
}

func (e *Emulator) filled(x, y int) bool {
	return e.frameBuffer[x+y*BaseWidth] == 1
}

func bits(b byte) [8]byte {
	res := [8]byte{}
	bit := byte(0b10000000)
	for i := 0; i < 8; i++ {
		res[i] = (b & bit) >> (7 - i)
		bit >>= 1
	}
	return res
}
