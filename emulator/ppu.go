package emulator

import (
	"sort"
)

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

// ByXPosition implements sort.Interface for []SpriteObject based on
// the xPosition field.
type ByXPosition []SpriteObject

func (a ByXPosition) Len() int           { return len(a) }
func (a ByXPosition) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByXPosition) Less(i, j int) bool { return a[i].xPosition < a[j].xPosition }

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

	// Window counter
	windowLineCounter uint8
}

func getBufferPositionFromXY(x, y, width int) int {
	return y*width + x
}

func (e *Emulator) getPPUMode() uint8 {
	ly := e.mem.GetLY()

	var mode uint8
	if ly >= 144 {
		mode = 1
	} else if ly < 144 {
		if e.ppu.ppuDot <= lcdMode2Bounds {
			mode = 2
		} else if e.ppu.ppuDot <= lcdMode3Bounds {
			mode = 3
		} else {
			mode = 0
		}
	}

	return mode
}

func (e *Emulator) isInterruptModeEnable(mode uint8) bool {
	status := e.mem.GetLCDStatus()

	switch mode {
	case 0:
		return status.Mode0HBlankSTATInterruptSource
	case 1:
		return status.Mode1VBlankSTATInterruptSource
	case 2:
		return status.Mode2OAMSTATInterruptSource
	}

	return false
}

func (e *Emulator) lyCompare() {
	if e.ppu.ppuDot != 0 {
		return
	}

	status := e.mem.GetLCDStatus()
	ly := e.mem.GetLY()
	lyc := e.mem.GetLYC()

	if ly == lyc {
		e.mem.lcdStatus.LYCLYFlag = true
		if status.LYCLYSTATInterruptSource {
			e.mem.requestInterruptLCDStat()
		}
	} else {
		e.mem.lcdStatus.LYCLYFlag = false
	}
}

// Set the status of the LCD based on the current state of memory.
func (e *Emulator) updateLCDStatus() {
	status := e.mem.GetLCDStatus()

	previousMode := status.ModeFlag
	status.ModeFlag = e.getPPUMode()

	if status.ModeFlag != previousMode {
		if status.ModeFlag == 0 {
			//gb.Memory.doHDMATransfer()
		}

		if e.isInterruptModeEnable(status.ModeFlag) {
			e.mem.requestInterruptLCDStat()
		}
	}
	e.lyCompare()

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

func (e *Emulator) isWindowVisible() bool {
	lcdc := e.mem.GetLCDC()
	wx := e.mem.GetWX()
	wy := e.mem.GetWY()
	ly := e.mem.GetLY()

	if wx <= 166 && wy <= 143 && lcdc.WindowEnable && ly > wy {
		return true
	}

	return false
}

func (e *Emulator) isInsideWindow(currentX uint8) bool {
	lcdc := e.mem.GetLCDC()
	ly := e.mem.GetLY()
	wx := e.mem.GetWX()
	wy := e.mem.GetWY()

	isWindow := false
	if lcdc.WindowEnable && ly >= wy && currentX+7 >= wx {
		isWindow = true
	}

	return isWindow
}

func (e *Emulator) getXYWindow(currentX uint8) (uint8, uint8) {
	wx := e.mem.GetWX()

	var x = currentX - wx + 7
	var y = e.ppu.windowLineCounter

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
	var tileIndex = sprite.tileIndex
	if e.getObjectHeight() == 16 {
		// Bit 0 of tile index for 8x16 objects should be ignored
		tileIndex = sprite.tileIndex & 0xFE
		if pixelY >= 8 {
			tileIndex = tileIndex | 0x01
			pixelY -= 8
		}
	}

	spriteColor := e.getColorFromVRAM(int(tileIndex), pixelX, pixelY)
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

	// Sort sprites by priority
	// In Non-CGB mode, the smaller the X coordinate, the higher the priority. When X coordinates are identical, the object located first in OAM has higher priority.
	// In CGB mode, only the object’s location in OAM determines its priority. The earlier the object, the higher its priority.
	sort.Stable(ByXPosition(e.ppu.spritesSelected[:e.ppu.numSpritesSelected]))
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
