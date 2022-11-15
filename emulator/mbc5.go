package emulator

// MBC5
// Max 256 KiB ROM and 512x4 bits RAM
// 4 "Register"
// - 0x0000-1FFF - RAM Enable (Write Only)
// - 0x2000-2FFF - 8 least significant bits of ROM bank number (Write Only)
// - 0x3000-3FFF - 9th bit of ROM bank number (Write Only)
// - 0x4000-5FFF - RAM bank number (Write Only)
type MBC5 struct {
	*BaseMBC

	romBank       uint16
	ramBank       uint8
	enableRamBank bool
	enableRumble  bool
}

func NewMBC5(baseMBC *BaseMBC) *MBC5 {
	mbc := &MBC5{
		BaseMBC: baseMBC,
		romBank: 1,
	}

	return mbc
}

func (mbc *MBC5) Write(addr uint16, val uint8) {
	// 4 "Register"
	// - 0x0000-1FFF - RAM Enable
	// - 0x2000-2FFF - 8 least significant bits of ROM bank number
	// - 0x3000-3FFF - 9th bit of ROM bank number
	// - 0x4000-5FFF - RAM bank number

	switch addr >> 13 {
	case 0: // 0x0000-1FFF - RAM Enable
		// Ram Enabled
		if val == 0x0A {
			mbc.enableRamBank = true
		} else {
			mbc.enableRamBank = false
		}
	case 1: // 0x2000-3FFF
		if addr <= 0x2FFF {
			//  0x2000-2FFF - 8 least significant bits of ROM bank number
			mbc.romBank = (mbc.romBank & 0x0100) | uint16(val)
			mbc.romBank &= uint16(mbc.RomBanks) - 1
		} else if addr <= 0x3FFF {
			// 0x4000 - 0x5FFF - 9th bit of ROM bank number
			mbc.romBank = (uint16(val&0x01) << 8) | (mbc.romBank & 0x00FF)
			mbc.romBank &= uint16(mbc.RomBanks) - 1
		}
	case 2: // 0x4000-5FFF - RAM bank number
		mbc.enableRumble = GetBit(val, 3)
		mbc.ramBank = val & 0x0F
		mbc.ramBank &= uint8(mbc.RamBanks) - 1
	case 5: // 0xA000 - 0xBFFF
		if mbc.enableRamBank {
			// RAM sizes are 8 KiB, 32 KiB and 128 KiB.
			ramAddr := (uint32(mbc.ramBank) << 13) | uint32(addr&0x1fff)
			mbc.ram[ramAddr] = val
		}
	}
}

func (mbc *MBC5) Read(addr uint16) uint8 {
	// It can map up to 64 Mbits (8 MiB) of ROM.

	switch addr >> 13 {
	case 0, 1: // 0x0000–0x3FFF
		// Contains the first 16 KiB of the ROM.
		return mbc.rom[addr&0x3fff]
	case 2, 3: // 0x4000 - 0x7FFF
		bankAddr := (uint32(mbc.romBank) << 14) + uint32(addr&0x3fff)
		return mbc.rom[bankAddr]
	case 5: // 0xA000 - 0xBFFF
		// RAM sizes are 8 KiB, 32 KiB and 128 KiB.
		if mbc.enableRamBank {
			ramAddr := (uint32(mbc.ramBank) << 13) | uint32(addr&0x1fff)
			return mbc.ram[ramAddr]
		}
		return 0xFF
	}

	return 0
}
