package emulator

func (e *Emulator) memoryBankController5Write(addr uint16, val uint8) {
	// 4 "Register"
	// - 0x0000-1FFF - RAM Enable
	// - 0x2000-2FFF - 8 least significant bits of ROM bank number
	// - 0x3000-3FFF - 9th bit of ROM bank number
	// - 0x4000-5FFF - RAM bank number

	switch addr >> 13 {
	case 0: // 0x0000-1FFF - RAM Enable
		// Ram Enabled
		if val == 0x0A {
			e.mbc5EnableRamBank = true
		} else {
			e.mbc5EnableRamBank = false
		}
	case 1: // 0x2000-3FFF
		if addr <= 0x2FFF {
			//  0x2000-2FFF - 8 least significant bits of ROM bank number
			e.mbc5RomBank = uint16(val) | (e.mbc5RomBank & 0x0100)
		} else if addr <= 0x3FFF {
			// 0x4000 - 0x5FFF - 9th bit of ROM bank number
			e.mbc5RomBank = uint16(val) | (e.mbc5RomBank & 0x0100)
		}
	case 2: // 0x4000-5FFF - RAM bank number
		e.mbc5EnableRumble = GetBit(val, 3)
		e.mbc5RamBank = val & 0x0F
	case 5: // 0xA000 - 0xBFFF
		if e.mbc5EnableRamBank {
			// RAM sizes are 8 KiB, 32 KiB and 128 KiB.
			ramAddr := uint32(e.mbc5RamBank)<<13 | uint32(addr&0x1fff)
			e.extrambank[ramAddr] = val
		}
	}
}

func (e *Emulator) memoryBankController5Read(addr uint16) uint8 {
	// It can map up to 64 Mbits (8 MiB) of ROM.

	switch addr >> 13 {
	case 0, 1: // 0x0000–0x3FFF
		// Contains the first 16 KiB of the ROM.
		return e.rom0[addr&0x3fff]
	case 2, 3: // 0x4000 - 0x7FFF
		bankAddr := uint32(e.mbc5RomBank)<<14 + uint32(addr&0x3fff)
		return e.rom0[bankAddr]
	case 5: // 0xA000 - 0xBFFF
		// RAM sizes are 8 KiB, 32 KiB and 128 KiB.
		if e.mbc5EnableRamBank {
			ramAddr := uint32(e.mbc5RamBank)<<13 | uint32(addr&0x1fff)
			return e.extrambank[ramAddr]
		}
		return 0xFF
	}

	return 0
}

func (e *Emulator) memoryBankController5(addr uint16, val uint8, write bool) uint8 {
	// Max 256 KiB ROM and 512x4 bits RAM
	// 4 "Register"
	// - 0x0000-1FFF - RAM Enable (Write Only)
	// - 0x2000-2FFF - 8 least significant bits of ROM bank number (Write Only)
	// - 0x3000-3FFF - 9th bit of ROM bank number (Write Only)
	// - 0x4000-5FFF - RAM bank number (Write Only)

	if write {
		e.memoryBankController5Write(addr, val)
	} else {
		return e.memoryBankController5Read(addr)
	}

	return 0
}
