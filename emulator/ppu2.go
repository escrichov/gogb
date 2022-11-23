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
		e.setLCDStatus()
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

func (e *Emulator) PPURunCycleEnabled() bool {
	ly := e.mem.GetLY()
	renderFrame := false

	// Modes 2 (OAM scan) & 3 (Drawing pixels)
	if ly < HEIGHT {
		if e.ppu.ppuDot == 80 {
			e.oamScan()
		} else if e.ppu.ppuDot == 369 {
			e.drawScanline()
		}
	}

	// Scanline (Every 456 PPU Dots)
	if e.ppu.ppuDot == 456 {
		// Interrupt VBlank and draw frame
		if ly == (HEIGHT - 1) {
			e.mem.requestInterruptVBlank()
			renderFrame = true
		}

		// Increment Line
		ly = (ly + 1) % 154
		e.mem.SetLY(ly)

		// Increment window internal line counter
		if ly == 0 {
			e.ppu.windowLineCounter = 0
		} else {
			if e.isWindowVisible() {
				e.ppu.windowLineCounter = (e.ppu.windowLineCounter + 1) % 154
			}
		}
		e.ppu.ppuDot = 0
	} else {
		e.ppu.ppuDot++
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
	}

	return renderFrame
}
