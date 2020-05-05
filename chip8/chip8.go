package chip8

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

type Emulator struct {
	display     *Display
	frameBuffer [BufferSize]byte
	memory      [MemorySize]byte
	vReg        [VRegisterSize]byte
	iReg        uint16
	delayTimer  byte
	soundTimer  byte
	pc          uint16
	stack       []uint16
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
	rand.Seed(time.Now().UnixNano())
	t0 := time.NewTicker(ClockSpeed)
	t1 := time.NewTicker(FrameRate)
	t2 := time.NewTicker(TimerSpeed)
	for {
		select {
		case <-t0.C:
			e.Cycle()
		case <-t1.C:
			e.display.Render(e.frameBuffer)
		case <-t2.C:
			if e.delayTimer > 0 {
				e.delayTimer--
			}
			if e.soundTimer > 0 {
				if e.soundTimer == 1 {
					e.display.Beep()
				}
				e.soundTimer--
			}
		}
	}
}

func (e *Emulator) Cycle() {
	// fetch -> decode -> execute
	key, pressed := e.display.Key()
	opcode := uint16(e.memory[e.pc])<<8 | uint16(e.memory[e.pc+1])
	e.execute(opcode, key, pressed)
}

func (e *Emulator) decode(opcode uint16) Instruction {
	kk := opcode & 0x00FF
	n := opcode & 0x000F

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
	case 0x2000:
		return CALL
	case 0x3000:
		return SEVxByte
	case 0x4000:
		return SNEVxByte
	case 0x5000:
		return SEVxVy
	case 0x6000:
		return LDVxByte
	case 0x7000:
		return ADDVxByte
	case 0x8000:
		switch n {
		case 0x0:
			return LDVxVy
		case 0x1:
			return OR
		case 0x2:
			return AND
		case 0x3:
			return XOR
		case 0x4:
			return ADDVxVy
		case 0x5:
			return SUB
		case 0x6:
			return SHR
		case 0x7:
			return SUBN
		case 0xE:
			return SHL
		default:
			return UNKNOWN
		}
	case 0x9000:
		return SNEVxVy
	case 0xA000:
		return LDIAddr
	case 0xB000:
		return JPV0Addr
	case 0xC000:
		return RND
	case 0xD000:
		return DRW
	case 0xE000:
		switch kk {
		case 0x9E:
			return SKP
		case 0xA1:
			return SKNP
		default:
			return UNKNOWN
		}
	case 0xF000:
		switch kk {
		case 0x07:
			return LDVxDT
		case 0x0A:
			return LDVxK
		case 0x15:
			return LDDTVx
		case 0x18:
			return LDSTVx
		case 0x1E:
			return ADDIVx
		case 0x29:
			return LDFVx
		case 0x33:
			return LDBVx
		case 0x55:
			return LDIVx
		case 0x65:
			return LDVxI
		default:
			return UNKNOWN
		}
	default:
		return UNKNOWN
	}
}

func (e *Emulator) execute(opcode uint16, key byte, pressed bool) {
	//e.descOpcode(opcode)

	inst := e.decode(opcode)
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	nnn := opcode & 0x0FFF
	kk := opcode & 0x00FF
	n := opcode & 0x000F
	incPC := true

	switch inst {
	case SYS:
		// 0nnn - SYS addr
		// Jump to a machine code routine at nnn.
		//
		// This instruction is only used on the old computers on which Chip-8
		// was originally implemented. It is ignored by modern interpreters.
		panic("not implemented: SYS")
	case CLS:
		// 00E0 - CLS
		// Clear the display.
		for i := 0; i < BufferSize; i++ {
			e.frameBuffer[i] = 0
		}
	case RET:
		// 00EE - RET
		// Return from a subroutine.
		//
		// The interpreter sets the program counter to the address at the top
		// of the stack, then subtracts 1 from the stack pointer.
		e.pc = e.stack[len(e.stack)-1]
		e.stack = e.stack[:len(e.stack)-1]
	case JPAddr:
		// 1nnn - JP addr
		// Jump to location nnn.
		//
		// The interpreter sets the program counter to nnn.
		e.pc = nnn
		incPC = false
	case CALL:
		// 2nnn - CALL addr
		// Call subroutine at nnn.
		//
		// The interpreter increments the stack pointer, then puts the current
		// PC on the top of the stack. The PC is then set to nnn.
		e.stack = append(e.stack, e.pc)
		e.pc = nnn
		incPC = false
	case SEVxByte:
		// 3xkk - SE Vx, byte
		// Skip next instruction if Vx = kk.
		//
		// The interpreter compares register Vx to kk, and if they are equal,
		// increments the program counter by 2.
		if uint16(e.vReg[x]) == kk {
			e.pc += 2
		}
	case SNEVxByte:
		// 4xkk - SNE Vx, byte
		// Skip next instruction if Vx != kk.
		//
		// The interpreter compares register Vx to kk, and if they are not
		// equal, increments the program counter by 2.
		if uint16(e.vReg[x]) != kk {
			e.pc += 2
		}
	case SEVxVy:
		// 5xy0 - SE Vx, Vy
		// Skip next instruction if Vx = Vy.
		//
		// The interpreter compares register Vx to register Vy, and if they are
		// equal, increments the program counter by 2.
		if e.vReg[x] == e.vReg[y] {
			e.pc += 2
		}
	case LDVxByte:
		// 6xkk - LD Vx, byte
		// Set Vx = kk.
		//
		// The interpreter puts the value kk into register Vx.
		e.vReg[x] = byte(kk)
	case ADDVxByte:
		// 7xkk - ADD Vx, byte
		// Set Vx = Vx + kk.
		//
		// Adds the value kk to the value of register Vx, then stores the
		// result in Vx.
		e.vReg[x] = byte((uint16(e.vReg[x]) + kk) & 0xFF)
	case LDVxVy:
		// 8xy0 - LD Vx, Vy
		// Set Vx = Vy.
		//
		// Stores the value of register Vy in register Vx.
		e.vReg[x] = e.vReg[y]
	case OR:
		// 8xy1 - OR Vx, Vy
		// Set Vx = Vx OR Vy.
		//
		// Performs a bitwise OR on the values of Vx and Vy, then stores the
		// result in Vx. A bitwise OR compares the corrseponding bits from two
		// values, and if either bit is 1, then the same bit in the result is
		// also 1. Otherwise, it is 0.
		e.vReg[x] |= e.vReg[y]
	case AND:
		// 8xy2 - AND Vx, Vy
		// Set Vx = Vx AND Vy.
		//
		// Performs a bitwise AND on the values of Vx and Vy, then stores the
		// result in Vx. A bitwise AND compares the corrseponding bits from two
		// values, and if both bits are 1, then the same bit in the result is
		// also 1. Otherwise, it is 0.
		e.vReg[x] &= e.vReg[y]
	case XOR:
		// 8xy3 - XOR Vx, Vy
		// Set Vx = Vx XOR Vy.
		//
		// Performs a bitwise exclusive OR on the values of Vx and Vy, then
		// stores the result in Vx. An exclusive OR compares the corresponding
		// bits from two values, and if the bits are not both the same, then
		//the corresponding bit in the result is set to 1. Otherwise, it is 0.
		e.vReg[x] ^= e.vReg[y]
	case ADDVxVy:
		// 8xy4 - ADD Vx, Vy
		// Set Vx = Vx + Vy, set VF = carry.
		//
		// The values of Vx and Vy are added together. If the result is greater
		// than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the
		// lowest 8 bits of the result are kept, and stored in Vx.
		added := uint16(e.vReg[x]) + uint16(e.vReg[y])
		if added > 0xFF {
			e.vReg[0xF] = 1
		} else {
			e.vReg[0xF] = 0
		}
		e.vReg[x] = byte(added)
	case SUB:
		// 8xy5 - SUB Vx, Vy
		// Set Vx = Vx - Vy, set VF = NOT borrow.
		//
		// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted
		// from Vx, and the results stored in Vx.
		if e.vReg[x] > e.vReg[y] {
			e.vReg[0xF] = 1
			e.vReg[x] -= e.vReg[y]
		} else if e.vReg[x] == e.vReg[y] {
			e.vReg[0xF] = 0
			e.vReg[x] = 0
		} else {
			e.vReg[0xF] = 0
			e.vReg[x] = 255 - (e.vReg[y] - e.vReg[x]) + 1
		}
	case SHR:
		// 8xy6 - SHR Vx {, Vy}
		// Set Vx = Vx SHR 1.
		//
		// If the least-significant bit of Vx is 1, then VF is set to 1,
		// otherwise 0. Then Vx is divided by 2.
		e.vReg[0xF] = e.vReg[x] & 0x1
		e.vReg[x] >>= 1
	case SUBN:
		// 8xy7 - SUBN Vx, Vy
		// Set Vx = Vy - Vx, set VF = NOT borrow.
		//
		// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted
		// from Vy, and the results stored in Vx.
		if e.vReg[y] > e.vReg[x] {
			e.vReg[0xF] = 1
			e.vReg[x] = e.vReg[y] - e.vReg[x]
		} else if e.vReg[y] == e.vReg[x] {
			e.vReg[0xF] = 0
			e.vReg[x] = 0
		} else {
			e.vReg[0xF] = 0
			e.vReg[x] = 255 - (e.vReg[x] - e.vReg[y]) + 1
		}
	case SHL:
		// 8xyE - SHL Vx {, Vy}
		// Set Vx = Vx SHL 1.
		//
		// If the most-significant bit of Vx is 1, then VF is set to 1,
		// otherwise to 0. Then Vx is multiplied by 2.
		e.vReg[0xF] = (e.vReg[x] & 0b10000000) >> 7
		e.vReg[x] <<= 1
	case SNEVxVy:
		// 9xy0 - SNE Vx, Vy
		// Skip next instruction if Vx != Vy.
		//
		// The values of Vx and Vy are compared, and if they are not equal, the
		// program counter is increased by 2.
		if e.vReg[x] != e.vReg[y] {
			e.pc += 2
		}
	case LDIAddr:
		// Annn - LD I, addr
		// Set I = nnn.
		//
		// The value of register I is set to nnn.
		e.iReg = nnn
	case JPV0Addr:
		// Bnnn - JP V0, addr
		// Jump to location nnn + V0.
		//
		// The program counter is set to nnn plus the value of V0.
		e.pc = nnn + uint16(e.vReg[0])
		incPC = false
	case RND:
		// Cxkk - RND Vx, byte
		// Set Vx = random byte AND kk.
		//
		// The interpreter generates a random number from 0 to 255, which is
		// then ANDed with the value kk. The results are stored in Vx. See
		// instruction 8xy2 for more information on AND.
		e.vReg[x] = byte(uint16(rand.Intn(256)) & kk)
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
	case SKP:
		// Ex9E - SKP Vx
		// Skip next instruction if key with the value of Vx is pressed.
		//
		// Checks the keyboard, and if the key corresponding to the value of Vx
		// is currently in the down position, PC is increased by 2.
		if pressed && key == e.vReg[x] {
			e.pc += 2
		}
	case SKNP:
		// ExA1 - SKNP Vx
		// Skip next instruction if key with the value of Vx is not pressed.
		//
		// Checks the keyboard, and if the key corresponding to the value of Vx
		// is currently in the up position, PC is increased by 2.
		if !pressed || key != e.vReg[x] {
			e.pc += 2
		}
	case LDVxDT:
		// Fx07 - LD Vx, DT
		// Set Vx = delay timer value.
		//
		// The value of DT is placed into Vx.
		e.vReg[x] = e.delayTimer
	case LDVxK:
		// Fx0A - LD Vx, K
		// Wait for a key press, store the value of the key in Vx.
		//
		// All execution stops until a key is pressed, then the value of that
		// key is stored in Vx.
		if pressed {
			e.vReg[x] = key
		} else {
			incPC = false
		}
	case LDDTVx:
		// Fx15 - LD DT, Vx
		// Set delay timer = Vx.
		//
		// DT is set equal to the value of Vx.
		e.delayTimer = e.vReg[x]
	case LDSTVx:
		// Fx18 - LD ST, Vx
		// Set sound timer = Vx.
		//
		// ST is set equal to the value of Vx.
		e.soundTimer = e.vReg[x]
	case ADDIVx:
		// Fx1E - ADD I, Vx
		// Set I = I + Vx.
		//
		// The values of I and Vx are added, and the results are stored in I.
		e.iReg += uint16(e.vReg[x])
	case LDFVx:
		// Set I = location of sprite for digit Vx.
		//
		// The value of I is set to the location for the hexadecimal sprite
		// corresponding to the value of Vx. See section 2.4, Display, for more
		// information on the Chip-8 hexadecimal font.
		e.iReg = uint16(e.vReg[x] * 5)
	case LDBVx:
		// Fx33 - LD B, Vx
		// Store BCD representation of Vx in memory locations I, I+1, and I+2.
		//
		// The interpreter takes the decimal value of Vx, and places the
		// hundreds digit in memory at location in I, the tens digit at
		// location I+1, and the ones digit at location I+2.
		hundreds := e.vReg[x] / 100
		tens := (e.vReg[x] / 10) % 10
		ones := e.vReg[x] % 10
		e.memory[e.iReg] = hundreds
		e.memory[e.iReg+1] = tens
		e.memory[e.iReg+2] = ones
	case LDIVx:
		// Fx55 - LD [I], Vx
		// Store registers V0 through Vx in memory starting at location I.
		//
		// The interpreter copies the values of registers V0 through Vx into
		// memory, starting at the address in I.
		for i := uint16(0); i < x+1; i++ {
			e.memory[e.iReg+i] = e.vReg[i]
		}
	case LDVxI:
		// Fx65 - LD Vx, [I]
		// Read registers V0 through Vx from memory starting at location I.
		//
		// The interpreter reads values from memory starting at location I into
		// registers V0 through Vx.
		for i := uint16(0); i < x+1; i++ {
			e.vReg[i] = e.memory[e.iReg+i]
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
	if x >= BaseWidth || y >= BaseHeight {
		return false
	}
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

func (e *Emulator) descOpcode(opcode uint16) {
	desc := fmt.Sprintf("0x%04x", opcode) + " "
	inst := e.decode(opcode)
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	nnn := opcode & 0x0FFF
	kk := opcode & 0x00FF
	n := opcode & 0x000F

	switch inst {
	case SYS:
		desc += fmt.Sprintf("SYS 0x%03x", nnn)
	case CLS:
		desc += "CLS"
	case RET:
		desc += "RET"
	case JPAddr:
		desc += fmt.Sprintf("JP 0x%03x", nnn)
	case CALL:
		desc += fmt.Sprintf("CALL 0x%03x", nnn)
	case SEVxByte:
		desc += fmt.Sprintf("SE V%x, 0x%02x", x, kk)
	case SNEVxByte:
		desc += fmt.Sprintf("SNE V%x, 0x%02x", x, kk)
	case SEVxVy:
		desc += fmt.Sprintf("SE V%x, V%x", x, y)
	case LDVxByte:
		desc += fmt.Sprintf("LD V%x, 0x%02x", x, kk)
	case ADDVxByte:
		desc += fmt.Sprintf("ADD V%x, 0x%02x", x, kk)
	case LDVxVy:
		desc += fmt.Sprintf("LD V%x, V%x", x, y)
	case OR:
		desc += fmt.Sprintf("OR V%x, V%x", x, y)
	case AND:
		desc += fmt.Sprintf("AND V%x, V%x", x, y)
	case XOR:
		desc += fmt.Sprintf("XOR V%x, V%x", x, y)
	case ADDVxVy:
		desc += fmt.Sprintf("ADD V%x, V%x", x, y)
	case SUB:
		desc += fmt.Sprintf("SUB V%x, V%x", x, y)
	case SHR:
		desc += fmt.Sprintf("SHR V%x, V%x", x, y)
	case SUBN:
		desc += fmt.Sprintf("SUBN V%x, V%x", x, y)
	case SHL:
		desc += fmt.Sprintf("SHL V%x, V%x", x, y)
	case SNEVxVy:
		desc += fmt.Sprintf("SNE V%x, V%x", x, y)
	case LDIAddr:
		desc += fmt.Sprintf("LD I, 0x%03x", nnn)
	case JPV0Addr:
		desc += fmt.Sprintf("JP V0, %03x", nnn)
	case RND:
		desc += fmt.Sprintf("RND V%x, byte", x)
	case DRW:
		desc += fmt.Sprintf("DRW V%x, V%x, 0x%x", x, y, n)
	case SKP:
		desc += fmt.Sprintf("SKP V%x", x)
	case SKNP:
		desc += fmt.Sprintf("SKNP V%x", x)
	case LDVxDT:
		desc += fmt.Sprintf("LD V%x, DT", x)
	case LDVxK:
		desc += fmt.Sprintf("LD V%x, K", x)
	case LDDTVx:
		desc += fmt.Sprintf("LD DT, V%x", x)
	case LDSTVx:
		desc += fmt.Sprintf("LD ST, V%x", x)
	case ADDIVx:
		desc += fmt.Sprintf("ADD I, V%x", x)
	case LDFVx:
		desc += fmt.Sprintf("LD F, V%x", x)
	case LDBVx:
		desc += fmt.Sprintf("LD B, V%x", x)
	case LDIVx:
		desc += fmt.Sprintf("LD [I], V%x", x)
	case LDVxI:
		desc += fmt.Sprintf("LD V%x, [I]", x)
	default:
		desc += "Unknown"
	}

	fmt.Println(desc)
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
