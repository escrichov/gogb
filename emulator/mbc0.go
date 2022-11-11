package emulator

func (e *Emulator) memoryBankController0(addr uint16, val uint8, write bool) uint8 {
	switch addr >> 13 {
	case 0, 1, 2, 3: // 0x0000 - 0x7FFF
		return e.rom0[addr]
	case 5: // 0xA000 - 0xBFFF
		return 0xFF
	}

	return 0
}
