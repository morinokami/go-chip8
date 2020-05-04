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

		e.execute(0x00E0)

		for _, b := range e.frameBuffer {
			if b != 0 {
				t.Fatal("Error: buffer not cleared")
			}
		}
	})

	t.Run("1nnn JPAddr", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)

		e.execute(0x1228)

		if e.pc != 0x228 {
			t.Errorf("want=0x%04x, got=0x%04x", 0x228, e.pc)
		}
	})

	t.Run("6xkk LDVxByte", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)

		e.execute(0x600C)

		if e.vReg[0] != 0x0C {
			t.Errorf("want=0x%04x, got=0x%04x", 0x0C, e.vReg[0])
		}
	})

	t.Run("7xkk ADDVxByte", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)

		e.execute(0x7009)

		if e.vReg[0] != 0x09 {
			t.Errorf("want=0x%04x, got=0x%04x", 0x09, e.vReg[0])
		}
	})

	t.Run("Annn LDIAddr", func(t *testing.T) {
		d := NewDisplay()
		e := New(d)

		e.execute(0xA22A)

		if e.iReg != 0x22A {
			t.Errorf("want=0x%04x, got=0x%04x", 0x22A, e.iReg)
		}
	})

	t.Run("Dxyn DRW", func(t *testing.T) {
		// pass
	})

}
