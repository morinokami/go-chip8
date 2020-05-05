package chip8

import (
	"math/rand"
	"testing"
	"time"
)

func TestEmulator(t *testing.T) {

	t.Run("00E0 CLS", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < BufferSize; i++ {
			e.frameBuffer[i] = byte(rand.Intn(2))
		}

		e.execute(0x00E0, 0, false)

		for _, b := range e.frameBuffer {
			if b != 0 {
				t.Fatal("Error: buffer not cleared")
			}
		}
	})

	t.Run("00EE RET", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.stack = append(e.stack, 0x666)

		e.execute(0x00EE, 0, false)

		if e.pc != 0x666+2 {
			t.Errorf("got=0x%04x, want=0x%04x", e.pc, 0x666+2)
		}
	})

	t.Run("1nnn JPAddr", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)

		e.execute(0x1228, 0, false)

		if e.pc != 0x228 {
			t.Errorf("got=0x%04x, want=0x%04x", e.pc, 0x228)
		}
	})

	t.Run("2nnn CALL", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		prevPC := e.pc

		e.execute(0x2242, 0, false)

		if e.pc != 0x242 {
			t.Errorf("got=0x%04x, want=0x%04x", e.pc, 0x242)
		}
		if e.stack[0] != prevPC {
			t.Errorf("got=0x%04x, want=0x%04x", e.stack[0], prevPC)
		}
	})

	t.Run("3xkk SEVxByte", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[2] = 1
		prevPC := e.pc

		e.execute(0x3201, 0, false)

		if e.pc != prevPC+4 {
			t.Errorf("got=0x%04x, want=0x%04x", e.pc, prevPC+4)
		}
	})

	t.Run("4xkk SNEVxByte", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		prevPC := e.pc

		e.execute(0x452A, 0, false)

		if e.pc != prevPC+4 {
			t.Errorf("got=0x%04x, want=0x%04x", e.pc, prevPC+4)
		}
	})

	t.Run("5xy0 SEVxVy", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[5] = 123
		e.vReg[6] = 123
		prevPC := e.pc

		e.execute(0x5560, 0, false)

		if e.pc != prevPC+4 {
			t.Errorf("got=0x%04x, want=0x%04x", e.pc, prevPC+4)
		}
	})

	t.Run("6xkk LDVxByte", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)

		e.execute(0x600C, 0, false)

		if e.vReg[0] != 0x0C {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 0x0C)
		}
	})

	t.Run("7xkk ADDVxByte", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)

		e.execute(0x7009, 0, false)

		if e.vReg[0] != 0x09 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 0x09)
		}
	})

	t.Run("8xy0 LDVxVy", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0xD] = 121
		e.vReg[0xE] = 123

		e.execute(0x8DE0, 0, false)

		if e.vReg[0xD] != e.vReg[0xE] {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0xD], e.vReg[0xE])
		}
	})

	t.Run("8xy1 OR", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 5
		e.vReg[1] = 2

		e.execute(0x8011, 0, false)

		if e.vReg[0] != 7 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 7)
		}
	})

	t.Run("8xy2 AND", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 5
		e.vReg[1] = 3

		e.execute(0x8012, 0, false)

		if e.vReg[0] != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 1)
		}
	})

	t.Run("8xy3 XOR", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 5
		e.vReg[1] = 3

		e.execute(0x8013, 0, false)

		if e.vReg[0] != 6 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 6)
		}
	})

	t.Run("8xy4 ADDVxVy", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 0xFF
		e.vReg[1] = 0x1

		e.execute(0x8014, 0, false)

		if e.vReg[0] != 0 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 0)
		}
		if e.vReg[0xF] != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0xF], 1)
		}
	})

	t.Run("8xy5 SUB", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 1
		e.vReg[1] = 2

		e.execute(0x8015, 0, false)

		if e.vReg[0] != 0xFF {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 0xFF)
		}
		if e.vReg[0xF] != 0 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0xF], 0)
		}
	})

	t.Run("8xy6 SHR", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 7

		e.execute(0x8006, 0, false)

		if e.vReg[0xF] != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0xF], 1)
		}
		if e.vReg[0] != 3 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 3)
		}
	})

	t.Run("8xy7 SUBN", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 1
		e.vReg[1] = 2

		e.execute(0x8017, 0, false)

		if e.vReg[0] != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 1)
		}
		if e.vReg[0xF] != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0xF], 1)
		}
	})

	t.Run("8xyE SHL", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 0xFF

		e.execute(0x800E, 0, false)

		if e.vReg[0xF] != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0xF], 1)
		}
		if e.vReg[0] != 0xFE {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 0xFE)
		}
	})

	t.Run("9xy0 SNEVxVy", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[5] = 121
		e.vReg[6] = 123
		prevPC := e.pc

		e.execute(0x9560, 0, false)

		if e.pc != prevPC+4 {
			t.Errorf("got=0x%04x, want=0x%04x", e.pc, prevPC+4)
		}
	})

	t.Run("Annn LDIAddr", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)

		e.execute(0xA22A, 0, false)

		if e.iReg != 0x22A {
			t.Errorf("got=0x%04x, want=0x%04x", e.iReg, 0x22A)
		}
	})

	t.Run("Bnnn JPV0Addr", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 1

		e.execute(0xB228, 0, false)

		if e.pc != 0x229 {
			t.Errorf("got=0x%04x, want=0x%04x", e.pc, 0x229)
		}
	})

	t.Run("Cxkk RND", func(t *testing.T) {
		// pass
	})

	t.Run("Dxyn DRW", func(t *testing.T) {
		// pass
	})

	t.Run("Ex9E SKP", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[9] = 3
		prevPC := e.pc

		e.execute(0xE99E, 0x3, true)

		if e.pc != prevPC+4 {
			t.Errorf("got=0x%04x, want=0x%04x", e.pc, prevPC+4)
		}
	})

	t.Run("ExA1 SKNP", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[9] = 3
		prevPC := e.pc

		e.execute(0xE9A1, 0x2, true)

		if e.pc != prevPC+4 {
			t.Errorf("got=0x%04x, want=0x%04x", e.pc, prevPC+4)
		}
	})

	t.Run("Fx07 LDVxDT", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.delayTimer = 1

		e.execute(0xF007, 0, false)

		if e.vReg[0] != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 1)
		}
	})

	t.Run("Fx0A LDVxK", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		prevPC := e.pc

		e.execute(0xF00A, 0, false)

		if e.pc != prevPC {
			t.Fatalf("got=0x%04x, want=0x%04x", e.pc, prevPC)
		}

		e.execute(0xF00A, 1, true)
		if e.vReg[0] != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.vReg[0], 1)
		}
	})

	t.Run("Fx15 LDDTVx", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 1

		e.execute(0xF015, 0, false)

		if e.delayTimer != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.delayTimer, 1)
		}
	})

	t.Run("Fx18 LDSTVx", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.vReg[0] = 1

		e.execute(0xF018, 0, false)

		if e.soundTimer != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.soundTimer, 1)
		}
	})

	t.Run("Fx1E ADDIVx", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.iReg = 7
		e.vReg[2] = 4

		e.execute(0xF21E, 0, false)

		if e.iReg != 11 {
			t.Errorf("got=0x%04x, want=0x%04x", e.iReg, 11)
		}
	})

	t.Run("Fx29 LDFVx", func(t *testing.T) {
		// pass
	})

	t.Run("Fx33 LDBVx", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.iReg = 1
		e.vReg[0] = 213

		e.execute(0xF033, 0, false)

		if e.memory[1] != 2 {
			t.Errorf("got=0x%04x, want=0x%04x", e.memory[1], 2)
		}
		if e.memory[2] != 1 {
			t.Errorf("got=0x%04x, want=0x%04x", e.memory[2], 1)
		}
		if e.memory[3] != 3 {
			t.Errorf("got=0x%04x, want=0x%04x", e.memory[3], 3)
		}
	})

	t.Run("Fx55 LDIVx", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.iReg = 1
		for i := byte(0); i < 9; i++ {
			e.vReg[i] = i
		}

		e.execute(0xF855, 0, false)

		for i := byte(0); i < 9; i++ {
			if e.memory[1+i] != i {
				t.Errorf("got=0x%04x, want=0x%04x", e.memory[1+i], i)
			}
		}
	})

	t.Run("Fx65 LDVxI", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)
		e.iReg = 1
		for i := byte(0); i < 9; i++ {
			e.memory[1+i] = i
		}

		e.execute(0xF865, 0, false)

		for i := byte(0); i < 9; i++ {
			if e.vReg[i] != i {
				t.Errorf("got=0x%04x, want=0x%04x", e.vReg[i], i)
			}
		}
	})

}
