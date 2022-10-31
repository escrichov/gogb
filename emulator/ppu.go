package emulator

type LCDControl struct {
	LCDPPUEnable           bool // Bit 7, 0=Off, 1=On
	WindowTileMapArea      bool // Bit 6, 0=9800-9BFF, 1=9C00-9FFF
	WindowEnable           bool // Bit 5, 0=Off, 1=On
	BgWindowTileDataArea   bool // Bit 4, 0=8800-97FF, 1=8000-8FFF
	BgTileMapArea          bool // Bit 3, 0=9800-9BFF, 1=9C00-9FFF
	ObjSize                bool // Bit 2, 0=8x8, 1=8x16
	ObjEnable              bool // Bit 1, 0=Off, 1=On
	BgWindowEnablePriority bool // Bit 0, 0=Off, 1=On
}

func (e *Emulator) GetLCDC() *LCDControl {
	return &e.lcdcControl
}

func (e *Emulator) SetLCDC(value uint8) {
	e.io[320] = value
	e.lcdcControl.LCDPPUEnable = GetBit(value, 7)
	e.lcdcControl.WindowTileMapArea = GetBit(value, 6)
	e.lcdcControl.WindowEnable = GetBit(value, 5)
	e.lcdcControl.BgWindowTileDataArea = GetBit(value, 4)
	e.lcdcControl.BgTileMapArea = GetBit(value, 3)
	e.lcdcControl.ObjSize = GetBit(value, 2)
	e.lcdcControl.ObjEnable = GetBit(value, 1)
	e.lcdcControl.BgWindowEnablePriority = GetBit(value, 0)
}

func (e *Emulator) getColor(tile, yOffset, xOffset int) uint8 {
	videoRamIndex := tile*16 + yOffset*2
	tileData := e.videoRam[videoRamIndex]
	tileData1 := e.videoRam[videoRamIndex+1]
	return ((tileData1>>xOffset)%2)*2 + (tileData>>xOffset)%2
}

func (e *Emulator) PPURun() {
	// PPU
	div := e.GetDIV()
	e.SetDIV(div + e.cycles - e.prevCycles)
	for ; e.prevCycles != e.cycles; e.prevCycles++ {
		lcdc := e.GetLCDC()
		if lcdc.LCDPPUEnable {
			e.ppuDot++

			// Render Scanline (Every 256 PPU Dots)
			if e.ppuDot == 456 {
				ly := e.GetLY()

				// Only render visible lines (up to line 144)
				if ly < HEIGHT {
					for tmp := WIDTH - 1; tmp >= 0; tmp-- {

						// IsWindow
						isWindow := false
						if lcdc.WindowEnable && ly >= e.io[330] && uint8(tmp) >= (e.io[331]-7) {
							isWindow = true
						}

						// xOffset
						var xOffset uint8
						if isWindow {
							xOffset = uint8(tmp) - e.io[331] + 7
						} else {
							xOffset = uint8(tmp) + e.io[323]
						}

						// yOffset
						var yOffset uint8
						if isWindow {
							yOffset = ly - e.io[330]
						} else {
							yOffset = ly + e.io[322]
						}

						// PaletteIndex
						var paletteIndex uint16 = 0

						// Tile
						tileMapArea := lcdc.BgTileMapArea
						if isWindow {
							tileMapArea = lcdc.WindowTileMapArea
						}

						videoRamIndex := uint16(6)
						if tileMapArea {
							videoRamIndex = 7
						}
						videoRamIndex = videoRamIndex<<10 | uint16(yOffset)/8*32 + uint16(xOffset)/8
						var tile = e.videoRam[videoRamIndex]

						// Color
						var tileValue int
						if lcdc.BgWindowTileDataArea {
							tileValue = int(tile)
						} else {
							tileValue = 256 + int(int8(tile))
						}
						color := e.getColor(tileValue, int(yOffset&7), int(7-xOffset&7))

						// Sprites
						if lcdc.ObjEnable {
							for spriteIndex := uint8(0); spriteIndex < WIDTH; spriteIndex += 4 {
								spriteX := uint8(tmp) - e.io[spriteIndex+1] + 8
								spriteY := ly - e.io[spriteIndex] + 16

								spriteYOffset := uint8(0)
								if (e.io[spriteIndex+3] & 64) != 0 {
									spriteYOffset = 7
								}
								spriteYOffset = spriteY ^ spriteYOffset

								spriteXOffset := uint8(7)
								if (e.io[spriteIndex+3] & 32) != 0 {
									spriteXOffset = 0
								}
								spriteXOffset = spriteX ^ spriteXOffset

								spriteColor := e.getColor(int(e.io[spriteIndex+2]), int(spriteYOffset), int(spriteXOffset))

								if spriteX < 8 && spriteY < 8 && !((e.io[spriteIndex+3]&128) != 0 && color != 0) && spriteColor != 0 {
									color = spriteColor
									if e.io[spriteIndex+3]&16 == 0 {
										paletteIndex = uint16(1)
									} else {
										paletteIndex = uint16(2)
									}
									break
								}
							}
						}

						paletteIndexValue := uint16((e.io[327+paletteIndex]>>(2*color))%4) + paletteIndex*4&7
						frameBufferIndex := uint16(ly)*WIDTH + uint16(tmp)
						e.frameBuffer[frameBufferIndex] = e.palette[paletteIndexValue]
					}
				}

				if ly == (HEIGHT - 1) {
					e.SetIF(e.GetIF() | 1)

					e.renderFrame()
					e.manageKeyboardEvents()
				}

				// Increment Line
				e.SetLY((ly + 1) % 154)
				e.ppuDot = 0
			}
		} else {
			e.SetLY(0)
			e.ppuDot = 0
		}
	}
}
