package ppu

import "emulator-go/emulator/gb/utils"

type PixelFetcher struct {
	fetcherHorizontalPosition int
	fetcherState              int
	fetcherTileNumber         byte
	fetcherDataLo             byte
	fetcherDataHi             byte
	fetcherWindowActive       bool
	windowVerticalPosition    int

	fetcherSpritesActive      bool
	fetcherSpritesState       int
	fetcherSpritesSpriteIndex int
	fetcherSpritesDataLow     byte
	fetcherSpritesDataHigh    byte
}

func (ppu *PPU) GetBackgroundWindowTileAddress() uint16 {
	var address int
	var row, col int
	if ppu.pixelFetcher.fetcherWindowActive {
		if ppu.lcdControl.WindowTileMapArea {
			address = 0x9C00
		} else {
			address = 0x9800
		}
		col = ppu.pixelFetcher.windowVerticalPosition
		row = ppu.pixelFetcher.windowVerticalPosition
	} else {
		if ppu.lcdControl.BgTileMapArea {
			address = 0x9C00
		} else {
			address = 0x9800
		}
		col = ((ppu.GetScrollColumn() / 8) + ppu.pixelFetcher.fetcherHorizontalPosition) & 0x1F
		row = (ppu.scanLine.lineIndex + ppu.GetScrollRow()) & 0xFF
	}

	return uint16(address + utils.RowColtoPos(row, col, 32))
}

func (ppu *PPU) GetTileDataAddress() uint16 {
	tileIndex := 0
	var backgroundTileAddress uint16
	var addressingMethod int
	if ppu.lcdControl.BgWindowTileDataArea {
		backgroundTileAddress = uint16(0x8000)
		addressingMethod = AddressingMethodUnsigned
	} else {
		backgroundTileAddress = uint16(0x9000)
		addressingMethod = AddressingMethodSigned
	}

	if addressingMethod == AddressingMethodSigned {
		backgroundTileAddress = uint16(int32(backgroundTileAddress) + int32(int8(tileIndex))*16)
	} else {
		backgroundTileAddress = backgroundTileAddress + uint16(tileIndex)*16
	}

	return backgroundTileAddress
}

func (ppu *PPU) PixelFetcherStepGetTile() {
	address := ppu.GetBackgroundWindowTileAddress()
	ppu.pixelFetcher.fetcherTileNumber = ppu.MMU.Memory[address]
}

func (ppu *PPU) PixelFetcherStepGetDataLow() {
	address := ppu.GetTileDataAddress()
	ppu.pixelFetcher.fetcherDataLo = ppu.MMU.Memory[address]
}

func (ppu *PPU) PixelFetcherStepGetDataHigh() {
	address := ppu.GetTileDataAddress() + 1
	ppu.pixelFetcher.fetcherDataHi = ppu.MMU.Memory[address]
}

func (ppu *PPU) PixelFetcherStepPush() {
}

func (ppu *PPU) PixelFetcherStepSleep() {
}

func (ppu *PPU) GetSpriteToFetch() (bool, int) {
	return false, 0
}

func (ppu *PPU) WindowStartFetch() bool {
	if ppu.lcdControl.WindowEnable && !ppu.pixelFetcher.fetcherWindowActive {
		return true
	}
	return false
}

func (ppu *PPU) BackgroundFetcherStep() {
	switch ppu.pixelFetcher.fetcherState {
	case 0:
	case 1:
		ppu.PixelFetcherStepGetTile()
	case 2:
	case 3:
		ppu.PixelFetcherStepGetDataLow()
	case 4:
	case 5:
		ppu.PixelFetcherStepGetDataHigh()
	case 6:
	case 7:
		ppu.PixelFetcherStepPush()
	}

	ppu.pixelFetcher.fetcherState++
}

func (ppu *PPU) SpriteFetcherStep() {
	switch ppu.pixelFetcher.fetcherSpritesState {
	case 0:
	case 1:
		ppu.PixelFetcherStepGetTile()
	case 2:
	case 3:
		ppu.PixelFetcherStepGetDataLow()
	case 4:
	case 5:
		ppu.PixelFetcherStepGetDataHigh()
	case 6:
	case 7:
		ppu.PixelFetcherStepPush()
	}

	ppu.pixelFetcher.fetcherSpritesState++
}
