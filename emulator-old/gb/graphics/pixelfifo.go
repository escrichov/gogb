package ppu

type Pixel struct {
	Color              uint8
	Palette            uint8
	SpritePriority     uint8
	BackgroundPriority bool
}

type PixelFIFO struct {
	queue []Pixel
}

func (fifo *PixelFIFO) Enqueue(pixel Pixel) {
	fifo.queue = append(fifo.queue, pixel)
}

func (fifo *PixelFIFO) Deque() *Pixel {
	if len(fifo.queue) > 0 {
		pixel := &fifo.queue[0]
		fifo.queue = fifo.queue[1:]
		return pixel
	}
	return nil
}

func (fifo *PixelFIFO) Clear() {
	fifo.queue = fifo.queue[:0]
}

func (ppu *PPU) PixelFifoGetPixel() uint32 {
	return 0
}
