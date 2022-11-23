package emulator

func (e *Emulator) drawScanline() {
	lcdc := e.mem.GetLCDC()
	ly := e.mem.GetLY()
	framebuffer := e.window.GetFramebuffer()

	for tmp := WIDTH - 1; tmp >= 0; tmp-- {
		currentX := uint8(tmp)
		var color PPUColor

		// Get Background or Window pixel
		if lcdc.BgWindowEnablePriority {
			color = e.getBGWindowColor(currentX)
		} else {
			color = PPUColor{colorIndex: 0, paletteIndex: paletteBGP}
		}

		// Get sprite pixel
		if lcdc.ObjEnable {
			// Traverse only selected sprites
			for i := 0; i < e.ppu.numSpritesSelected; i++ {
				// If sprite has priority override color
				sprite := &e.ppu.spritesSelected[i]
				if e.isSpriteIsInColumn(sprite, currentX) {
					spriteColor := e.getSpriteColor(sprite, currentX)
					if e.spriteHasPriorityOverBG(sprite, spriteColor.colorIndex, color.colorIndex) {
						color = spriteColor
						break
					}
				}

			}
		}

		// Save color in framebuffer
		paletteColor := e.getColorFromPalette(&color)
		frameBufferIndex := getBufferPositionFromXY(int(currentX), int(ly), WIDTH)
		framebuffer[frameBufferIndex] = e.getDisplayColor(color.paletteIndex, paletteColor)
	}
}

func (e *Emulator) PPURun() bool {
	renderFrame := false

	// PPU
	cyclesElapsed := e.cycles - e.prevCycles
	for i := uint64(0); i < cyclesElapsed; i++ {
		renderFrame = e.PPURunCycle()
		if renderFrame {
			e.window.renderFrame()
			if e.showWindow {
				e.manageKeyboardEvents()
			}
		}
	}

	return renderFrame
}

func (ppu *PPU) isEndOfScanline() bool {
	if ppu.ppuDot == 0 {
		return true
	} else {
		return false
	}
}

func (ppu *PPU) isBeginOfVBlank(ly uint8) bool {
	if ly == HEIGHT {
		return true
	} else {
		return false
	}
}

func (e *Emulator) incrementWindowLineCounter(ly uint8) {
	if ly == 0 {
		e.ppu.windowLineCounter = 0
	} else {
		if e.isWindowVisible() {
			e.ppu.windowLineCounter++
		}
	}
}

func (e *Emulator) PPURunCycleEnabled() bool {
	ly := e.mem.GetLY()
	renderFrame := false

	// Scanline (Every 456 PPU Dots)
	e.ppu.ppuDot = (e.ppu.ppuDot + 1) % 456
	if e.ppu.isEndOfScanline() {
		// Increment Line
		ly = (ly + 1) % 154
		e.mem.SetLY(ly)

		// Interrupt VBlank and draw frame
		if e.ppu.isBeginOfVBlank(ly) {
			e.mem.requestInterruptVBlank()
			renderFrame = true
		}

		// Increment window internal line counter
		e.incrementWindowLineCounter(ly)
	}

	e.updateLCDStatus()

	// Modes 2 (OAM scan) & 3 (Drawing pixels)
	if ly < HEIGHT {
		if e.ppu.ppuDot == 80 {
			e.oamScan()
		} else if e.ppu.ppuDot == 369 {
			e.drawScanline()
		}
	}

	return renderFrame
}

func (e *Emulator) PPURunCycle() bool {
	renderFrame := false

	lcdc := e.mem.GetLCDC()
	if lcdc.LCDPPUEnable {
		renderFrame = e.PPURunCycleEnabled()
	} else {
		e.mem.SetLY(0)
		e.ppu.windowLineCounter = 0
		e.ppu.ppuDot = 0

		status := e.mem.GetLCDStatus()
		status.LYCLYFlag = false
		status.ModeFlag = 0
		e.mem.SaveLCDStatus()
	}

	return renderFrame
}
