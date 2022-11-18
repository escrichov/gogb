package emulator

const (
	lcdMode2Bounds = 80
	lcdMode3Bounds = lcdMode2Bounds + 172
)

const (
	paletteBGP  = uint8(0)
	paletteOBP0 = 1
	paletteOBP1 = 2
)

type PPU struct {
	paletteBGP []int32
	paletteOB0 []int32
	paletteOB1 []int32
	ppuDot     int
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

func (e *Emulator) getColorFromVRAM(tileNumber, xPixel int, yPixel int) uint8 {
	tileSize := 16   // 16 Bytes per tile (8x8 pixels of 2 bits each)
	bytesPerRow := 2 // 2 Bytes per row
	videoRamIndex := tileNumber*tileSize + yPixel*bytesPerRow
	tileData := e.mem.videoRam[videoRamIndex]    // Least Significant bits (LSB)
	tileData1 := e.mem.videoRam[videoRamIndex+1] // Most Significant bits (MSB)

	bitLSB := (tileData >> xPixel) & 0x01
	bitMSB := (tileData1 >> xPixel) & 0x01

	return bitMSB<<1 + bitLSB
}

func (e *Emulator) isInsideWindow(currentX uint8) bool {
	lcdc := e.mem.GetLCDC()
	ly := e.mem.GetLY()
	wx := e.mem.GetWX()
	wy := e.mem.GetWY()

	isWindow := false
	if lcdc.WindowEnable && ly >= wy && currentX >= (wx-7) {
		isWindow = true
	}

	return isWindow
}

func (e *Emulator) getXYWindow(currentX uint8) (uint8, uint8) {
	ly := e.mem.GetLY()
	wx := e.mem.GetWX()
	wy := e.mem.GetWY()

	var x = currentX - wx + 7
	var y = ly - wy

	return x, y
}

func (e *Emulator) getXYBackground(currentX uint8) (uint8, uint8) {
	ly := e.mem.GetLY()
	scx := e.mem.GetSCX()
	scy := e.mem.GetSCY()

	var x = currentX + scx
	var y = ly + scy

	return x, y
}

func (e *Emulator) getXYBackgroundWindow(currentX uint8, isWindow bool) (uint8, uint8) {
	if isWindow {
		return e.getXYWindow(currentX)
	} else {
		return e.getXYBackground(currentX)
	}
}

func (e *Emulator) getTilemapArea(isWindow bool) uint16 {
	lcdc := e.mem.GetLCDC()

	// Get Tile map area
	tileMapArea := lcdc.BgTileMapArea
	if isWindow {
		tileMapArea = lcdc.WindowTileMapArea
	}

	videoRamIndexPrefix := uint16(0x1800) // Tilemap $9800
	if tileMapArea {
		videoRamIndexPrefix = 0x1C00 // Tilemap $9C00
	}

	return videoRamIndexPrefix
}

func (e *Emulator) getTileNumber(isWindow bool, x, y uint8) uint8 {
	// Get Tile map area
	tileMapArea := e.getTilemapArea(isWindow)

	// 8x8 Tile Size
	tileSizeX := uint16(8)
	tileSizeY := uint16(8)

	// Tilemap size 32 x 32 => In each row there are 32 tiles
	tilesPerRow := uint16(32)

	tileX := uint16(x) / tileSizeX
	tileY := uint16(y) / tileSizeY * tilesPerRow
	videoRamIndex := tileMapArea + tileY + tileX
	return e.mem.videoRam[videoRamIndex]
}

func (e *Emulator) getTileIndexFromTileNumber(tileNumber uint8) int {
	lcdc := e.mem.GetLCDC()

	if lcdc.BgWindowTileDataArea {
		return int(tileNumber) // 0x8800 - 0x8FFF, 0x8000 + 0, 256
	} else {
		return 256 + int(int8(tileNumber)) // 0x8800 - 0x97FF, 8000 + (128, 383)
	}
}

func (e *Emulator) getBGWindowColor(currentX uint8) uint8 {
	// IsWindow
	isWindow := e.isInsideWindow(currentX)

	// xOffset
	xOffset, yOffset := e.getXYBackgroundWindow(currentX, isWindow)

	// Get tile number from tilemap
	var tileNumber = e.getTileNumber(isWindow, xOffset, yOffset)

	// Tile Index
	var tileIndex = e.getTileIndexFromTileNumber(tileNumber)

	// Color
	yPixel := int(yOffset & 7)
	xPixel := 7 - int(xOffset&7) // TODO: Why is reversed??
	color := e.getColorFromVRAM(tileIndex, xPixel, yPixel)

	return color
}

func (e *Emulator) spriteHasPriorityOverBG(spriteAddress uint8, spriteColor uint8, backgroundColor uint8) bool {
	oamFlags := e.mem.io[spriteAddress+3]
	// OAM
	// Bit7   BG and Window over OBJ (0=No, 1=BG and Window colors 1-3 over the OBJ)

	// Has priority if all conditions are met:
	// * Sprite color is not transparent (color != 0)
	// * OAM Flag bit Bit7 (BG and Window over OBJ) == 0 || OAM Flag bit Bit7 (BG and Window over OBJ) == 1 and backgroundColor = 0
	if spriteColor == 0 {
		return false
	}

	if GetBit(oamFlags, 7) {
		if backgroundColor == 0 {
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

func (e *Emulator) getSpriteXY(spriteAddress uint8, currentX uint8) (uint8, uint8) {
	oamYPosition := e.mem.io[spriteAddress]
	oamXPosition := e.mem.io[spriteAddress+1]
	oamFlags := e.mem.io[spriteAddress+3]
	ly := e.mem.GetLY()

	spriteX := currentX - oamXPosition + 8
	spriteY := ly - oamYPosition + 16

	// Check y flip
	spriteYOffset := uint8(0)
	if GetBit(oamFlags, 6) {
		spriteYOffset = 7
	}
	spriteYOffset = spriteY ^ spriteYOffset

	// Check x flip
	spriteXOffset := uint8(7)
	if GetBit(oamFlags, 5) {
		spriteXOffset = 0
	}
	spriteXOffset = spriteX ^ spriteXOffset

	return spriteXOffset, spriteYOffset
}

func (e *Emulator) getSpriteIfHasPriority(spriteIndex uint8, currentX uint8, backgroundColor uint8) (bool, uint8, uint8) {
	spriteAddress := spriteIndex * 4
	oamTileIndex := e.mem.io[spriteAddress+2]
	oamFlags := e.mem.io[spriteAddress+3]

	spriteX, spriteY := e.getSpriteXY(spriteAddress, currentX)
	spriteColor := e.getColorFromVRAM(int(oamTileIndex), int(spriteX), int(spriteY))

	if spriteX < 8 && spriteY < 8 && e.spriteHasPriorityOverBG(spriteAddress, spriteColor, backgroundColor) {
		var paletteIndex uint8
		if GetBit(oamFlags, 7) {
			paletteIndex = paletteOBP1 // OBP1
		} else {
			paletteIndex = paletteOBP0 // OBP0
		}
		return true, spriteColor, paletteIndex
	}

	return false, 0, 0
}

func (e *Emulator) getPaletteValue(paletteIndex uint8) uint8 {
	return e.mem.io[327+uint16(paletteIndex)]
}

func (e *Emulator) getColorFromPalette(paletteIndex uint8, colorIndex uint8) uint8 {
	paletteValue := e.getPaletteValue(paletteIndex)
	color := (paletteValue >> (colorIndex * 2)) & 0x3
	return color
}

func (e *Emulator) getDisplayColor(paletteIndex uint8, gbColor uint8) int32 {
	switch paletteIndex {
	case paletteBGP:
		return e.ppu.paletteBGP[gbColor]
	case paletteOBP0:
		return e.ppu.paletteOB0[gbColor]
	case paletteOBP1:
		return e.ppu.paletteOB1[gbColor]
	default:
		return e.ppu.paletteBGP[gbColor]
	}
}

func (e *Emulator) proccessScanline() {
	lcdc := e.mem.GetLCDC()
	ly := e.mem.GetLY()
	framebuffer := e.window.GetFramebuffer()

	for tmp := WIDTH - 1; tmp >= 0; tmp-- {
		currentX := uint8(tmp)

		colorIndex := e.getBGWindowColor(currentX)
		paletteIndex := paletteBGP

		// Sprites
		if lcdc.ObjEnable {
			// Traverse all sprites
			for spriteIndex := uint8(0); spriteIndex < WIDTH; spriteIndex++ {
				// If sprite has priority override color
				hasPriority, spriteColor, spritePalette := e.getSpriteIfHasPriority(spriteIndex, currentX, colorIndex)
				if hasPriority {
					colorIndex = spriteColor
					paletteIndex = spritePalette
				}
			}
		}

		color := e.getColorFromPalette(paletteIndex, colorIndex)
		frameBufferIndex := uint16(ly)*WIDTH + uint16(currentX)
		framebuffer[frameBufferIndex] = e.getDisplayColor(paletteIndex, color)
	}
}

func (e *Emulator) PPURun() bool {
	renderFrame := false
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
					e.proccessScanline()
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
