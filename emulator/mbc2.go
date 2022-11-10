package emulator

func (e *Emulator) memoryBankController2Write(addr uint16, val uint8) {
	// 1 "Register"
	// - 0x0000-0x3FFF - RAM Enable, ROM Bank Number

	switch addr >> 13 {
	case 0, 1: // RAM Enable, ROM Bank Number
		if GetBit16(addr, 8) {
			// ROM bank
			e.mbc2RomBank = val & 0x0f
			if e.mbc2RomBank == 0 {
				e.mbc2RomBank = 1
			}
			e.mbc2RomBank = e.mbc2RomBank & (uint8(e.romHeader.RomBanks) - 1)
		} else {
			// Ram Enabled
			if val&0x0f == 0x0A {
				e.mbc2EnableRamBank = true
			} else {
				e.mbc2EnableRamBank = false
			}
		}
	case 5: // 0xA000 - 0xBFFF
		if e.mbc2EnableRamBank {
			ramAddr := addr & 0x1ff
			if addr <= 0xA1FF {
				e.extrambank[ramAddr] = (val & 0x0F) | 0xF0
			} else if addr <= 0xBFFF {
				e.extrambank[ramAddr] = (val & 0x0F) | 0xF0
			}
		}
	}
}

func (e *Emulator) memoryBankController2Read(addr uint16) uint8 {
	// In its default configuration, MBC1 supports up to 512 KiB ROM with up to 32 KiB of banked RAM.

	switch addr >> 13 {
	case 0, 1: // 0x0000–0x3FFF
		// Contains the first 16 KiB of the ROM.
		return e.rom0[addr&0x3fff]
	case 2, 3: // 0x4000 - 0x7FFF
		bankAddr := uint32(e.mbc2RomBank)<<14 + uint32(addr&0x3fff)
		return e.rom0[bankAddr]
	case 5: // 0xA000 - 0xBFFF
		// 512 half-bytes of RAM
		if e.mbc2EnableRamBank {
			ramAddr := addr & 0x1ff
			if addr <= 0xA1FF {
				return e.extrambank[ramAddr]
			} else if addr <= 0xBFFF {
				return e.extrambank[ramAddr]
			} else {
				return 0xFF
			}
		}
		return 0xFF
	}

	return 0
}

func (e *Emulator) memoryBankController2(addr uint16, val uint8, write bool) uint8 {
	// Max 256 KiB ROM and 512x4 bits RAM
	// 1 "Register"
	// - 0x0000-0x3FFF - RAM Enable, ROM Bank Number [write-only]

	if write {
		e.memoryBankController2Write(addr, val)
	} else {
		return e.memoryBankController2Read(addr)
	}

	return 0
}
