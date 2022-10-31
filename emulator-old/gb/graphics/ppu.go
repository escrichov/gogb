package ppu

import (
	mmu "emulator-go/emulator/gb/memory"
	"emulator-go/emulator/gb/utils"
	"fmt"
)

type PixelFIFOMode struct {
	Mode                  int
	Cycles                int
	NextModeIndex         int
	VRAMAccessible        bool
	CGBPalettesAccessible bool
	OAMAccessible         bool
}

const (
	PixelFIFOMode0OAMScan    int = 0
	PixelFIFOMode1VBlank         = 1
	PixelFIFOMode2HBlank         = 2
	PixelFIFOMode3DrawPixels     = 3
)

var PixelFIFOModes = [4]PixelFIFOMode{
	PixelFIFOMode{PixelFIFOMode0OAMScan, 204, 2, true, true, true},
	PixelFIFOMode{PixelFIFOMode1VBlank, 4560, 2, true, true, true},
	PixelFIFOMode{PixelFIFOMode2HBlank, 80, 3, true, true, false},
	PixelFIFOMode{PixelFIFOMode3DrawPixels, 289, 0, false, false, false},
}

// TODO: PPU
// UpdateGraphics
//   * isLCDEnabled (Util)
//   * setLCDStatus
//     - clearScreen
//     - isLCDEnabled (Util)
//     - drawScanline
//       - renderSprites
//         - setPixel (Util)
//         - getColour
//       - renderTiles
//         - getTileSettings
//         - setTilePixel
//           - setPixel (Util)
//           - getColour

type Line struct {
	modeStep            int
	lineIndex           int
	cycle               int
	NumOamEntries       int
	OamEntries          [10]OAMEntry
	BackgroundPixelFIFO PixelFIFO
	SpritesPixelFIFO    PixelFIFO
}

type PPU struct {
	MMU               *mmu.MMU
	BackgroundScreen  []uint32
	Screen            []uint32
	CurrentModeIndex  int
	scanLine          Line
	windowLineCounter int
	lcdControl        LCDControl
	pixelFetcher      PixelFetcher
}

const (
	AddressingMethodUnsigned int = 0
	AddressingMethodSigned       = 1
)

type LCDControl struct {
	LCDPPUEnable           bool // 0=Off, 1=On
	WindowTileMapArea      bool // 0=9800-9BFF, 1=9C00-9FFF
	WindowEnable           bool // 0=Off, 1=On
	BgWindowTileDataArea   bool // 0=8800-97FF, 1=8000-8FFF
	BgTileMapArea          bool // 0=9800-9BFF, 1=9C00-9FFF
	ObjSize                bool // 0=8x8, 1=8x16
	ObjEnable              bool // 0=Off, 1=On
	BgWindowEnablePriority bool // 0=Off, 1=On
}

func (lcdControl *LCDControl) GetObjHeight() int {
	if lcdControl.ObjSize {
		return 16
	} else {
		return 8
	}
}

func (ppu *PPU) Init(MMU *mmu.MMU) {
	ppu.MMU = MMU
	ppu.BackgroundScreen = make([]uint32, 256*256)
	ppu.Screen = make([]uint32, 160*144)
}

func (ppu *PPU) GetScrollRow() int {
	return int(ppu.MMU.Read(0xFF43))
}

func (ppu *PPU) GetScrollColumn() int {
	return int(ppu.MMU.Read(0xFF42))
}

func (ppu *PPU) GetWindowRow() byte {
	return ppu.MMU.Read(0xFF4A)
}

func (ppu *PPU) GetWindowColumn() byte {
	return ppu.MMU.Read(0xFF4B)
}

func (ppu *PPU) ReadLCDControl() {

	lcdControlRegister := ppu.MMU.Read(0xFF40)

	ppu.lcdControl.LCDPPUEnable = utils.GetBit(lcdControlRegister, 7)
	ppu.lcdControl.WindowTileMapArea = utils.GetBit(lcdControlRegister, 6)
	ppu.lcdControl.WindowEnable = utils.GetBit(lcdControlRegister, 5)
	ppu.lcdControl.BgWindowTileDataArea = utils.GetBit(lcdControlRegister, 4)
	ppu.lcdControl.BgTileMapArea = utils.GetBit(lcdControlRegister, 3)
	ppu.lcdControl.ObjSize = utils.GetBit(lcdControlRegister, 2)
	ppu.lcdControl.ObjEnable = utils.GetBit(lcdControlRegister, 1)
	ppu.lcdControl.BgWindowEnablePriority = utils.GetBit(lcdControlRegister, 0)
}

func (ppu *PPU) drawTile(tile []byte, tileIndex int) {
	tileRow := (tileIndex / 32) * 8
	tileCol := (tileIndex % 32) * 8
	for i := 0; i < len(tile); i += 2 {
		left := tile[i]
		right := tile[i+1]
		for j := 0; j < 8; j++ {
			bitPos := 7 - j
			leftBit := utils.GetBit(left, bitPos)
			rightBit := utils.GetBit(right, bitPos)

			screenRow := tileRow + (i / 2)
			screenCol := tileCol + j
			screenPosition := utils.RowColtoPos(screenRow, screenCol, 256)
			if leftBit && rightBit {
				ppu.BackgroundScreen[screenPosition] = 0xffff0000
			} else if !leftBit && rightBit {
				ppu.BackgroundScreen[screenPosition] = 0xff00ff00
			} else if leftBit && !rightBit {
				ppu.BackgroundScreen[screenPosition] = 0xff777777
			} else {
				ppu.BackgroundScreen[screenPosition] = 0xffffffff
			}

		}
	}
}

func (ppu *PPU) drawTiles(backgroundTileMapAddress, backgroundTileAddress uint16, addressingMethod int) {
	for i := uint16(0); i < 32*32; i++ {
		tileIndex := ppu.MMU.Memory[backgroundTileMapAddress+i]
		var address uint16
		if addressingMethod == AddressingMethodSigned {
			address = uint16(int32(backgroundTileAddress) + int32(int8(tileIndex))*16)
		} else {
			address = backgroundTileAddress + uint16(tileIndex)*16
		}
		ppu.drawTile(ppu.MMU.Memory[address:address+16], int(i))
	}
}

func (ppu *PPU) Update() {

	ppu.ReadLCDControl()

	topScrollRow := ppu.GetScrollColumn()
	topScrollColumn := ppu.GetScrollRow()
	fmt.Println(topScrollRow, topScrollColumn)

	if !ppu.lcdControl.LCDPPUEnable {
		ppu.ClearScreen(0xffffffff)
		return
	}

	ppu.RunModes()

	//backgroundTileAddress := uint16(0x9000)
	//addressingMethod := AddressingMethodSigned
	//if lcdControl.BgWindowTileDataArea {
	//	backgroundTileAddress = uint16(0x8000)
	//	addressingMethod = AddressingMethodUnsigned
	//}
	//
	//// 0=9800-9BFF, 1=9C00-9FFF
	//backgroundTileMapAddress := uint16(0x9800)
	//if lcdControl.BgTileMapArea {
	//	backgroundTileMapAddress = uint16(0x9C00)
	//}
	//
	//ppu.drawTiles(backgroundTileMapAddress, backgroundTileAddress, addressingMethod)
	//
	//topScrollRow = 0
	//topScrollColumn = 0
	//ppu.drawScrollView(250, 0)

	//ppu.drawBackgroundRectangle(topScrollRow, topScrollColumn, 160, 144)
}

func (ppu *PPU) RunModes() {
	mode := PixelFIFOModes[ppu.CurrentModeIndex]

	switch mode.Mode {
	case PixelFIFOMode0OAMScan:
		ppu.PPUMode0()
	case PixelFIFOMode1VBlank:
		ppu.PPUMode1()
	case PixelFIFOMode2HBlank:
		ppu.PPUMode2()
	case PixelFIFOMode3DrawPixels:
		ppu.PPUMode3()
	}

	ppu.scanLine.cycle++
}

func (ppu *PPU) drawScrollView(topScrollRow, topScrollColumn int) {
	for column := 0; column < 160; column++ {
		for row := 0; row < 144; row++ {
			rowBackground := (topScrollRow + row) % 256
			columnBackground := (topScrollColumn + column) % 256
			ppu.Screen[utils.RowColtoPos(row, column, 160)] = ppu.BackgroundScreen[utils.RowColtoPos(rowBackground, columnBackground, 256)]
		}
	}
}

func (ppu *PPU) drawBackgroundRectangle(topRow int, topColumn int, width int, height int) {

	// Horizontal lines
	for c := 0; c < width; c++ {
		// Initial line
		rowPosition := topRow
		columnPosition := (topColumn + c) % 256
		ppu.BackgroundScreen[utils.RowColtoPos(rowPosition, columnPosition, 256)] = 0xFF00FF00

		// End line
		rowPosition = (rowPosition + height) % 256
		ppu.BackgroundScreen[utils.RowColtoPos(rowPosition, columnPosition, 256)] = 0xFF00FF00
	}

	// Vertical lines
	for r := 0; r < height; r++ {
		// Initial line
		rowPosition := (topRow + r) % 256
		columnPosition := topColumn
		ppu.BackgroundScreen[utils.RowColtoPos(rowPosition, columnPosition, 256)] = 0xFF00FF00

		// End line
		columnPosition = (topColumn + width) % 256
		ppu.BackgroundScreen[utils.RowColtoPos(rowPosition, columnPosition, 256)] = 0xFF00FF00
	}

	// Draw top position
	ppu.BackgroundScreen[utils.RowColtoPos(topRow, topColumn, 256)] = 0xFF0000FF
}

func (ppu *PPU) ClearScreen(color uint32) {
	for i, _ := range ppu.BackgroundScreen {
		ppu.BackgroundScreen[i] = uint32(color)
	}
}

// PPUMode0 - HBlank Mode
func (ppu *PPU) PPUMode0() {
	if ppu.scanLine.cycle == 456 {
		if ppu.scanLine.lineIndex == 144 {
			ppu.CurrentModeIndex = PixelFIFOMode1VBlank
		} else {
			ppu.CurrentModeIndex = PixelFIFOMode0OAMScan
		}
		ppu.scanLine.cycle = 0
		ppu.scanLine.lineIndex++
	}
}

// PPUMode1 - VBlank Mode
func (ppu *PPU) PPUMode1() {
	if ppu.scanLine.cycle == 456 {
		ppu.scanLine.cycle = 0
		if ppu.scanLine.lineIndex == 153 {
			ppu.scanLine.lineIndex = 0
			ppu.CurrentModeIndex = PixelFIFOMode0OAMScan
		} else {
			ppu.scanLine.lineIndex++
		}
	}
}

func (ppu *PPU) PPUMode2() {
	if ppu.scanLine.cycle == 80 {
		ppu.CurrentModeIndex = PixelFIFOMode3DrawPixels
	}

	switch ppu.scanLine.modeStep {
	case 0:
		ppu.OAMStepParse()
		ppu.scanLine.modeStep = 1
	case 1:
		ppu.OAMStepCheck()
		ppu.scanLine.modeStep = 0
	}
}

func (ppu *PPU) PPUMode3() {
	ok, spriteIndex := ppu.GetSpriteToFetch()
	if ok {
		fmt.Println(spriteIndex)
	}

	if ppu.pixelFetcher.fetcherSpritesActive {
		ppu.SpriteFetcherStep()
	}

	ppu.BackgroundFetcherStep()

	if ppu.WindowStartFetch() {
		ppu.pixelFetcher.fetcherWindowActive = true
		ppu.pixelFetcher.fetcherState = 0
		ppu.pixelFetcher.fetcherHorizontalPosition = 0
		ppu.scanLine.BackgroundPixelFIFO.Clear()
	}

	pixel := ppu.PixelFifoGetPixel()
	fmt.Println(pixel)
}
