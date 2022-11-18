package emulator

type MBC2 struct {
	*BaseMBC

	romBank       uint8
	enableRamBank bool
}

// NewMBC2
// Max 256 KiB ROM and 512x4 bits RAM
// 1 "Register"
// - 0x0000-0x3FFF - RAM Enable, ROM Bank Number [write-only]
func NewMBC2(baseMBC *BaseMBC) *MBC2 {
	mbc := &MBC2{BaseMBC: baseMBC, romBank: 1}
	mbc.ram = make([]byte, 512) // 512 Fixed ram
	return mbc
}

func (mbc *MBC2) Write(addr uint16, val uint8) {
	// 1 "Register"
	// - 0x0000-0x3FFF - RAM Enable, ROM Bank Number

	switch addr >> 13 {
	case 0, 1: // RAM Enable, ROM Bank Number
		if GetBit16(addr, 8) {
			// ROM bank
			mbc.romBank = val & 0x0f
			if mbc.romBank == 0 {
				mbc.romBank = 1
			}
			mbc.romBank = mbc.romBank & (uint8(mbc.RomBanks) - 1)
		} else {
			// Ram Enabled
			if val&0x0f == 0x0A {
				mbc.enableRamBank = true
			} else {
				mbc.enableRamBank = false
			}
		}
	case 5: // 0xA000 - 0xBFFF
		if mbc.enableRamBank {
			ramAddr := addr & 0x1ff
			if addr <= 0xA1FF {
				mbc.ram[ramAddr] = (val & 0x0F) | 0xF0
			} else if addr <= 0xBFFF {
				mbc.ram[ramAddr] = (val & 0x0F) | 0xF0
			}
		}
	}
}

func (mbc *MBC2) Read(addr uint16) uint8 {
	// In its default configuration, MBC1 supports up to 512 KiB ROM with up to 32 KiB of banked RAM.

	switch addr >> 13 {
	case 0, 1: // 0x0000–0x3FFF
		// Contains the first 16 KiB of the ROM.
		return mbc.rom[addr&0x3fff]
	case 2, 3: // 0x4000 - 0x7FFF
		bankAddr := uint32(mbc.romBank)<<14 + uint32(addr&0x3fff)
		return mbc.rom[bankAddr]
	case 5: // 0xA000 - 0xBFFF
		// 512 half-bytes of RAM
		if mbc.enableRamBank {
			ramAddr := addr & 0x1ff
			if addr <= 0xA1FF {
				return mbc.ram[ramAddr]
			} else if addr <= 0xBFFF {
				return mbc.ram[ramAddr]
			} else {
				return 0xFF
			}
		}
		return 0xFF
	}

	return 0
}
