package emulator

import "testing"

func PixelsEqual(p1 *Pixel, p2 *Pixel) bool {
	if p1.color == p2.color && p1.backgroundPriority == p2.backgroundPriority &&
		p1.palette == p2.palette && p1.spritePriority == p2.spritePriority {
		return true
	} else {
		return false
	}
}

func TestPushPop(t *testing.T) {
	var fifo PixelFIFO

	fifo.Push(0x12, 0x34, 0, false)
	if fifo.numPixels != 8 {
		t.Fatalf(`Invalid number of pixels: %v != 8`, fifo.numPixels)
	}
	fifo.Push(0x56, 0x78, 1, true)
	if fifo.numPixels != 16 {
		t.Fatalf(`Invalid number of pixels: %v != 16`, fifo.numPixels)
	}

	expectedPixels := []Pixel{
		{color: 0, palette: 0, backgroundPriority: false},
		{color: 0, palette: 0, backgroundPriority: false},
		{color: 1, palette: 0, backgroundPriority: false},
		{color: 3, palette: 0, backgroundPriority: false},
		{color: 0, palette: 0, backgroundPriority: false},
		{color: 1, palette: 0, backgroundPriority: false},
		{color: 2, palette: 0, backgroundPriority: false},
		{color: 0, palette: 0, backgroundPriority: false},
		{color: 0, palette: 1, backgroundPriority: true},
		{color: 3, palette: 1, backgroundPriority: true},
		{color: 1, palette: 1, backgroundPriority: true},
		{color: 3, palette: 1, backgroundPriority: true},
		{color: 1, palette: 1, backgroundPriority: true},
		{color: 2, palette: 1, backgroundPriority: true},
		{color: 2, palette: 1, backgroundPriority: true},
		{color: 0, palette: 1, backgroundPriority: true},
	}

	for i := 0; i < 16; i++ {
		pixel, _ := fifo.Pop()

		if fifo.numPixels != 16-i-1 {
			t.Fatalf(`Invalid number of pixels: %v != %v`, fifo.numPixels, 16-i)
		}

		if !PixelsEqual(pixel, &expectedPixels[i]) {
			t.Fatalf(`Pixel %v not equal`, i)
		}
	}

	pushed := fifo.Push(0x9A, 0xBC, 1, false)
	if !pushed {
		t.Fatalf(`Expected value %v, got %v`, true, false)
	}

	expectedPixels2 := []Pixel{
		{color: 3, palette: 1, backgroundPriority: false},
		{color: 0, palette: 1, backgroundPriority: false},
		{color: 1, palette: 1, backgroundPriority: false},
		{color: 3, palette: 1, backgroundPriority: false},
		{color: 3, palette: 1, backgroundPriority: false},
		{color: 1, palette: 1, backgroundPriority: false},
		{color: 2, palette: 1, backgroundPriority: false},
		{color: 0, palette: 1, backgroundPriority: false},
	}

	for i := 0; i < 8; i++ {
		pixel, _ := fifo.Pop()

		if fifo.numPixels != 8-i-1 {
			t.Fatalf(`Invalid number of pixels: %v != %v`, fifo.numPixels, 8-i)
		}

		if !PixelsEqual(pixel, &expectedPixels2[i]) {
			t.Fatalf(`Pixel %v not equal`, i)
		}
	}

}

func TestPopEmpty(t *testing.T) {
	var fifo PixelFIFO

	_, popped := fifo.Pop()

	if popped {
		t.Fatalf(`Expected value %v, got %v`, false, true)
	}
}

func TestCanPushFull(t *testing.T) {
	var fifo PixelFIFO
	fifo.Push(0x00, 0x00, 0x00, true)
	fifo.Push(0x00, 0x00, 0x00, true)

	if fifo.CanPush() {
		t.Fatalf(`Expected value %v, got %v`, true, false)
	}
}

func TestCanPushEmpty(t *testing.T) {
	var fifo PixelFIFO

	if !fifo.CanPush() {
		t.Fatalf(`Expected value %v, got %v`, false, true)
	}
}

func TestCanPushHalfEmpty(t *testing.T) {
	var fifo PixelFIFO
	fifo.Push(0x00, 0x00, 0x00, true)

	if !fifo.CanPush() {
		t.Fatalf(`Expected value %v, got %v`, false, true)
	}
}

func TestCanPushFullPop7(t *testing.T) {
	var fifo PixelFIFO
	fifo.Push(0x00, 0x00, 0x00, true)
	fifo.Push(0x00, 0x00, 0x00, true)
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()

	if fifo.CanPush() {
		t.Fatalf(`Expected value %v, got %v`, false, true)
	}
}

func TestCanPushFullPop8(t *testing.T) {
	var fifo PixelFIFO
	fifo.Push(0x00, 0x00, 0x00, true)
	fifo.Push(0x00, 0x00, 0x00, true)
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()
	fifo.Pop()

	if !fifo.CanPush() {
		t.Fatalf(`Expected value %v, got %v`, true, false)
	}
}

func TestPushFull(t *testing.T) {
	var fifo PixelFIFO
	fifo.Push(0x00, 0x00, 0x00, true)
	fifo.Push(0x00, 0x00, 0x00, true)
	pushed := fifo.Push(0x00, 0x00, 0x00, true)

	if pushed {
		t.Fatalf(`Expected value %v, got %v`, false, true)
	}
}

func TestPush(t *testing.T) {
	var fifo PixelFIFO
	fifo.Push(0x00, 0x00, 0x00, true)
	pushed := fifo.Push(0x00, 0x00, 0x00, true)

	if !pushed {
		t.Fatalf(`Expected value %v, got %v`, true, false)
	}
}
