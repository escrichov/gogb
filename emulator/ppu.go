package emulator

const (
	lcdMode2Bounds = 80
	lcdMode3Bounds = lcdMode2Bounds + 172
)

type PPU struct {
	palette []int32
	ppuDot  int
}

func (e *Emulator) getColor(tile, yOffset, xOffset int) uint8 {
	videoRamIndex := tile*16 + yOffset*2
	tileData := e.mem.videoRam[videoRamIndex]
	tileData1 := e.mem.videoRam[videoRamIndex+1]
	return ((tileData1>>xOffset)%2)*2 + (tileData>>xOffset)%2
}

// Set the status of the LCD based on the current state of memory.
func (e *Emulator) setLCDStatus() {
	status := e.mem.GetLCDStatus()
	lcdcControl := e.mem.GetLCDC()

	if !lcdcControl.LCDPPUEnable {
		status.LYCLYSTATInterruptSource = false
		status.Mode2OAMSTATInterruptSource = false
		status.Mode1VBlankSTATInterruptSource = false
		status.Mode0HBlankSTATInterruptSource = false
		status.LYCLYFlag = false
		status.ModeFlag = 0
	}

	ly := e.mem.GetLY()
	currentMode := status.ModeFlag

	var mode uint8
	requestInterrupt := false

	switch {
	case ly >= 144:
		mode = 1
		requestInterrupt = status.Mode1VBlankSTATInterruptSource
	case e.ppu.ppuDot < lcdMode2Bounds:
		mode = 2
		requestInterrupt = status.Mode2OAMSTATInterruptSource
	case e.ppu.ppuDot < lcdMode3Bounds:
		mode = 3
	default:
		mode = 0
		requestInterrupt = status.Mode0HBlankSTATInterruptSource
		if mode != currentMode {
			//gb.Memory.doHDMATransfer()
		}
	}

	if requestInterrupt && mode != currentMode {
		e.mem.requestInterruptLCDStat()
	}

	// Check if LYC == LY (coincidence flag)
	lyc := e.mem.GetLYC()
	if ly == lyc {
		e.mem.lcdStatus.LYCLYFlag = true
		if e.mem.lcdStatus.LYCLYSTATInterruptSource {
			e.mem.requestInterruptLCDStat()
		}
	} else {
		e.mem.lcdStatus.LYCLYFlag = false
	}

	e.mem.SaveLCDStatus()
}

func (e *Emulator) PPURun() bool {
	renderFrame := false
	framebuffer := e.window.GetFramebuffer()
	e.setLCDStatus()

	// PPU
	cyclesElapsed := e.cycles - e.prevCycles
	for i := uint64(0); i < cyclesElapsed; i++ {
		lcdc := e.mem.GetLCDC()
		if lcdc.LCDPPUEnable {
			e.ppu.ppuDot++

			// Render Scanline (Every 256 PPU Dots)
			if e.ppu.ppuDot == 456 {
				ly := e.mem.GetLY()

				// Only render visible lines (up to line 144)
				if ly < HEIGHT {
					for tmp := WIDTH - 1; tmp >= 0; tmp-- {
						currentX := uint8(tmp)
						wy := e.mem.io[330]
						wx := e.mem.io[331]
						scy := e.mem.io[322]
						scx := e.mem.io[323]

						// IsWindow
						isWindow := false
						if lcdc.WindowEnable && ly >= wy && currentX >= (wx-7) {
							isWindow = true
						}

						// xOffset
						var xOffset uint8
						if isWindow {
							xOffset = currentX - wx + 7
						} else {
							xOffset = currentX + scx
						}

						// yOffset
						var yOffset uint8
						if isWindow {
							yOffset = ly - wy
						} else {
							yOffset = ly + scy
						}

						// PaletteIndex
						var paletteIndex uint16 = 0

						// Tile
						tileMapArea := lcdc.BgTileMapArea
						if isWindow {
							tileMapArea = lcdc.WindowTileMapArea
						}

						videoRamIndexPrefix := uint16(0x1800)
						if tileMapArea {
							videoRamIndexPrefix = 0x1c00
						}
						videoRamIndex := videoRamIndexPrefix | uint16(yOffset)/8*32 + uint16(xOffset)/8
						var tile = e.mem.videoRam[videoRamIndex]

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
								spriteX := currentX - e.mem.io[spriteIndex+1] + 8
								spriteY := ly - e.mem.io[spriteIndex] + 16

								spriteYOffset := uint8(0)
								// Check y flip
								if (e.mem.io[spriteIndex+3] & 64) != 0 {
									spriteYOffset = 7
								}
								spriteYOffset = spriteY ^ spriteYOffset

								spriteXOffset := uint8(7)
								// Check x flip
								if (e.mem.io[spriteIndex+3] & 32) != 0 {
									spriteXOffset = 0
								}
								spriteXOffset = spriteX ^ spriteXOffset

								spriteColor := e.getColor(int(e.mem.io[spriteIndex+2]), int(spriteYOffset), int(spriteXOffset))

								if spriteX < 8 && spriteY < 8 && !((e.mem.io[spriteIndex+3]&128) != 0 && color != 0) && spriteColor != 0 {
									color = spriteColor
									if e.mem.io[spriteIndex+3]&16 == 0 {
										paletteIndex = uint16(1)
									} else {
										paletteIndex = uint16(2)
									}
									break
								}
							}
						}

						paletteIndexValue := uint16((e.mem.io[327+paletteIndex]>>(2*color))%4) + paletteIndex*4&7
						frameBufferIndex := uint16(ly)*WIDTH + uint16(currentX)
						framebuffer[frameBufferIndex] = e.ppu.palette[paletteIndexValue]
					}
				}

				if ly == (HEIGHT - 1) {
					e.mem.requestInterruptVBlank()
					renderFrame = true
				}

				// Increment Line
				e.mem.SetLY((ly + 1) % 154)
				e.ppu.ppuDot = 0
			}
		} else {
			e.mem.SetLY(0)
			e.ppu.ppuDot = 0
		}
	}

	return renderFrame
}
