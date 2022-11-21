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

const (
	numSpritesInOAM   = uint8(40)
	spriteOAMSize     = uint8(4)
	tileSizeBytesVRAM = 16 // 16 Bytes per tile (8x8 pixels of 2 bits each)
	tileSizeX         = 8  // 8x8 Tile Size
	tileSizeY         = 8
	tilesPerRow       = 32 // Tilemap size 32 x 32 => In each row there are 32 tiles
)

type SpriteObject struct {
	tileIndex uint8
	xPosition uint8
	yPosition uint8
	flags     uint8

	bgWindowOverOBJ  bool
	yFlip            bool
	xFlip            bool
	paletteIndex     uint8
	tileVRAMBank     uint8
	paletteNumberCGB uint8
}

type PPUColor struct {
	paletteIndex uint8
	colorIndex   uint8
}

type PPU struct {
	paletteBGP []uint32
	paletteOB0 []uint32
	paletteOB1 []uint32
	ppuDot     int

	// OAM Scan
	spritesSelected    [10]SpriteObject
	numSpritesSelected int
}

func getBufferPositionFromXY(x, y, width int) int {
	return y*width + x
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

func (e *Emulator) getColorFromVRAM(tileNumber int, xPixel, yPixel uint8) uint8 {
	bytesPerRow := 2 // 2 Bytes per row
	videoRamIndex := tileNumber*tileSizeBytesVRAM + int(yPixel)*bytesPerRow
	tileData := e.mem.videoRam[videoRamIndex]    // Least Significant bits (LSB)
	tileData1 := e.mem.videoRam[videoRamIndex+1] // Most Significant bits (MSB)

	xShift := 7 - xPixel // Bit 7 represents the leftmost pixel, and bit 0 the rightmost.
	bitLSB := (tileData >> xShift) & 0x01
	bitMSB := (tileData1 >> xShift) & 0x01

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

func (e *Emulator) getXYSprite(sprite *SpriteObject, currentX uint8) (uint8, uint8) {
	ly := e.mem.GetLY()
	var x = currentX - sprite.xPosition + 8
	var y = ly - sprite.yPosition + 16

	// Flip x
	if sprite.xFlip {
		x = e.flipTileXPosition(x)
	}

	// Flip y
	if sprite.yFlip {
		y = e.flipTileYPosition(y)
	}

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

	tileX := int(x) / tileSizeX
	tileY := int(y) / tileSizeY
	videoRamIndex := tileMapArea + uint16(getBufferPositionFromXY(tileX, tileY, tilesPerRow))
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

func (e *Emulator) getTileIndex(isWindow bool, x, y uint8) int {
	// Get tile number from tilemap
	tileNumber := e.getTileNumber(isWindow, x, y)

	// Tile Index
	return e.getTileIndexFromTileNumber(tileNumber)
}

func (e *Emulator) getBGWindowColor(currentX uint8) PPUColor {
	// IsWindow
	isWindow := e.isInsideWindow(currentX)

	// xOffset
	xOffset, yOffset := e.getXYBackgroundWindow(currentX, isWindow)

	// Tile Index
	var tileIndex = e.getTileIndex(isWindow, xOffset, yOffset)

	// Color
	yPixel := yOffset & 7
	xPixel := xOffset & 7
	color := e.getColorFromVRAM(tileIndex, xPixel, yPixel)

	return PPUColor{paletteIndex: paletteBGP, colorIndex: color}
}

func (e *Emulator) spriteHasPriorityOverBG(sprite *SpriteObject, spriteColor uint8, backgroundColor uint8) bool {
	// Has priority if all conditions are met:
	// * Sprite color is not transparent (color != 0)
	if spriteColor == 0 {
		return false
	}

	// OAM
	// Bit7   BG and Window over OBJ (0=No, 1=BG and Window colors 1-3 over the OBJ)
	// * OAM Flag bit Bit7 (BG and Window over OBJ) == 0 || OAM Flag bit Bit7 (BG and Window over OBJ) == 1 and backgroundColor = 0
	if sprite.bgWindowOverOBJ {
		if backgroundColor == 0 {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

func (e *Emulator) flipTileXPosition(position uint8) uint8 {
	return position ^ 0x7 // This reverse values: 0-7 -> 7-0
}

func (e *Emulator) flipTileYPosition(position uint8) uint8 {
	height := e.getObjectHeight()
	return position ^ (height - 1) // This reverse values: 0-7 -> 7-0 or 0-15 -> 15-0
}

func (e *Emulator) getSpriteColor(sprite *SpriteObject, currentX uint8) PPUColor {
	pixelX, pixelY := e.getXYSprite(sprite, currentX)
	spriteColor := e.getColorFromVRAM(int(sprite.tileIndex), pixelX, pixelY)
	return PPUColor{paletteIndex: sprite.paletteIndex, colorIndex: spriteColor}
}

func (e *Emulator) getPaletteValue(paletteIndex uint8) uint8 {
	// BGP=327, OBP0=328, OBP1=329
	return e.mem.io[327+uint16(paletteIndex)]
}

func (e *Emulator) getColorFromPalette(ppuColor *PPUColor) uint8 {
	paletteValue := e.getPaletteValue(ppuColor.paletteIndex)
	color := (paletteValue >> (ppuColor.colorIndex * 2)) & 0x3
	return color
}

func (e *Emulator) getDisplayColor(paletteIndex uint8, gbColor uint8) uint32 {
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

func (e *Emulator) getObjectHeight() uint8 {
	lcdc := e.mem.GetLCDC()

	// Sprite height
	height := uint8(8)
	if lcdc.ObjSize {
		height = 16
	}

	return height
}

func (object *SpriteObject) saveSpriteData(spriteData []byte) {
	flags := spriteData[3]

	object.yPosition = spriteData[0]
	object.xPosition = spriteData[1]
	object.tileIndex = spriteData[2]
	object.flags = flags
	object.bgWindowOverOBJ = GetBit(flags, 7)
	object.yFlip = GetBit(flags, 6)
	object.xFlip = GetBit(flags, 5)
	if GetBit(flags, 4) {
		object.paletteIndex = paletteOBP1 // OBP1
	} else {
		object.paletteIndex = paletteOBP0 // OBP0
	}

	object.tileVRAMBank = BoolToUint8(GetBit(flags, 3))
}

func (e *Emulator) oamScan() {
	ly := e.mem.GetLY()
	e.ppu.numSpritesSelected = 0

	for spriteIndex := uint8(0); spriteIndex < numSpritesInOAM; spriteIndex++ {
		if e.isSpriteIsInLine(spriteIndex, ly) {
			spriteAddress := spriteIndex * spriteOAMSize
			e.ppu.spritesSelected[e.ppu.numSpritesSelected].saveSpriteData(e.mem.io[spriteAddress : spriteAddress+4])
			e.ppu.numSpritesSelected++
			if e.ppu.numSpritesSelected == 10 {
				break
			}
		}
	}

	// Sort sprites by x position
}

func (e *Emulator) isSpriteIsInLine(spriteIndex uint8, currentY uint8) bool {
	spriteAddress := spriteIndex * spriteOAMSize
	objectHeight := int(e.getObjectHeight())

	spriteYPosition := int(e.mem.io[spriteAddress]) - 16
	spriteLowerBound := spriteYPosition
	spriteUpperBound := spriteYPosition + objectHeight

	// Check if current line is between sprite bounds
	if int(currentY) >= spriteLowerBound && int(currentY) < spriteUpperBound {
		return true
	} else {
		return false
	}
}

func (e *Emulator) isSpriteIsInColumn(sprite *SpriteObject, currentX uint8) bool {
	spriteXPosition := int(sprite.xPosition) - 8
	spriteLowerBound := spriteXPosition
	spriteUpperBound := spriteXPosition + 8

	// Check if current line is between sprite bounds
	if int(currentX) >= spriteLowerBound && int(currentX) < spriteUpperBound {
		return true
	} else {
		return false
	}
}

func (e *Emulator) proccessScanline() {
	lcdc := e.mem.GetLCDC()
	ly := e.mem.GetLY()
	framebuffer := e.window.GetFramebuffer()

	e.oamScan()

	for tmp := WIDTH - 1; tmp >= 0; tmp-- {
		currentX := uint8(tmp)

		// Get Background or Window pixel
		color := e.getBGWindowColor(currentX)

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
