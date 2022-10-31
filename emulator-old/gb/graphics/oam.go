package ppu

import "emulator-go/emulator/gb/utils"

type OAMEntry struct {
	VerticalPosition   int
	HorizontalPosition int
	TileIndex          int
	BGWindowAOverObj   bool
	VerticalFlip       bool
	HorizontalFlip     bool
	PaletteNumber      int
	TileVRAMBank       int
	CGBPaletteNumber   int
}

func (oamEntry *OAMEntry) OAMSelect(currentHorizontalLine int, tileHeight int) bool {
	beginVerticalPosition := oamEntry.VerticalPosition - 16
	endVerticalPosition := beginVerticalPosition + tileHeight

	if beginVerticalPosition >= currentHorizontalLine && currentHorizontalLine < endVerticalPosition && oamEntry.HorizontalPosition > 0 {
		return true
	}

	return false
}

func (oamEntry *OAMEntry) Parse(data []byte) {
	oamEntry.HorizontalPosition = int(data[0])
	oamEntry.VerticalPosition = int(data[1])
	oamEntry.TileIndex = int(data[2])
	oamEntry.BGWindowAOverObj = utils.GetBit(data[2], 7)
	oamEntry.VerticalFlip = utils.GetBit(data[2], 6)
	oamEntry.HorizontalFlip = utils.GetBit(data[2], 5)
	oamEntry.PaletteNumber = utils.GetBitInt(data[2], 4)
	oamEntry.TileVRAMBank = utils.GetBitInt(data[3], 3)
	oamEntry.CGBPaletteNumber = int(data[3] & 0x07)
}

func (ppu *PPU) OAMStepParse() {
	if ppu.scanLine.NumOamEntries == 10 {
		return
	}
	address := uint16(0xFE00 + ppu.scanLine.cycle/2*4)
	ppu.scanLine.OamEntries[ppu.scanLine.NumOamEntries].Parse(ppu.MMU.Memory[address : address+4])
}

func (ppu *PPU) OAMStepCheck() {
	if ppu.scanLine.NumOamEntries == 10 {
		return
	}
	if ppu.scanLine.OamEntries[ppu.scanLine.NumOamEntries].OAMSelect(ppu.scanLine.lineIndex, ppu.lcdControl.GetObjHeight()) {
		ppu.scanLine.NumOamEntries++
	}
}
