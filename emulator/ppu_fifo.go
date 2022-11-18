package emulator

type Pixel struct {
	color              uint8 // value between 0 and 3
	palette            uint8 // on CGB a value between 0 and 7 and on DMG this only applies to sprites
	spritePriority     uint8 // on CGB this is the OAM index for the sprite and on DMG this doesn’t exist
	backgroundPriority bool  // holds the value of the OBJ-to-BG Priority bit
}

type PixelFIFO struct {
	pixels    [16]Pixel
	numPixels int
	pushIndex int
	popIndex  int
}

func (fifo *PixelFIFO) CanPush() bool {
	if fifo.numPixels <= 8 {
		return true
	} else {
		return false
	}
}

func (fifo *PixelFIFO) IsEmpty() bool {
	if fifo.numPixels == 0 {
		return true
	} else {
		return false
	}
}

func (fifo *PixelFIFO) Push(high uint8, low uint8, palette uint8, backgroundPriority bool) bool {
	if !fifo.CanPush() {
		return false
	}

	for i := 7; i >= 0; i-- {
		index := (fifo.pushIndex + i) % 16

		// Pixel decoding
		bitHigh := high & 0x01
		bitLow := low & 0x01
		color := (bitHigh << 1) | bitLow

		fifo.pixels[index].color = color
		fifo.pixels[index].palette = palette
		fifo.pixels[index].backgroundPriority = backgroundPriority

		high = high >> 1
		low = low >> 1
	}

	fifo.pushIndex = (fifo.pushIndex + 8) % 16
	fifo.numPixels += 8

	return true
}

func (fifo *PixelFIFO) Pop() (*Pixel, bool) {
	if fifo.IsEmpty() {
		return nil, false
	}

	pixel := &fifo.pixels[fifo.popIndex]

	fifo.popIndex = (fifo.popIndex + 1) % 16
	fifo.numPixels--

	return pixel, true
}

func (fifo *PixelFIFO) Clear() {
	fifo.numPixels = 0
}
