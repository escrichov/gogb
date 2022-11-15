package emulator

// MBC1
// In its default configuration, MBC1 supports up to 512 KiB ROM with up to 32 KiB of banked RAM.
// 4 "Registers"
// - 0x0000-0x1FFF - RAMG - MBC1 RAM gate register
// - 0x2000-0x3FFF - BANK1 - MBC1 bank register 1
// - 0x4000-0x5FFF - BANK2 - MBC1 bank register 2
// - 0x6000-0x7FFF - MODE - MBC1 mode register
type MBC1 struct {
	*BaseMBC
	bank1           uint8
	bank2           uint8
	memoryModel     int
	enableRamBank   bool
	allowedRomBank2 bool
	allowedRamBank2 bool
	isMBC1M         bool
}

func NewMBC1(baseMBC *BaseMBC) *MBC1 {
	mbc := &MBC1{bank1: 1, BaseMBC: baseMBC}
	if mbc.RomBanks >= 64 {
		mbc.allowedRomBank2 = true
	}
	if mbc.RamBanks >= 4 {
		mbc.allowedRamBank2 = true
	}

	return mbc
}

func (mbc *MBC1) memoryBankController1RomByBankNumber(bankNumber uint8, addr uint16) uint8 {
	bankAddr := uint32(bankNumber)<<14 + uint32(addr&0x3fff)
	return mbc.rom[bankAddr]
}

func (mbc *MBC1) memoryBankController1RomAddress00003FFF(addr uint16) uint8 {
	bankNumber := uint8(0)
	if mbc.memoryModel == 1 {
		if mbc.allowedRomBank2 {
			bankNumber = mbc.bank2 << 5
		}
	}
	return mbc.memoryBankController1RomByBankNumber(bankNumber, addr)
}

func (mbc *MBC1) memoryBankController1Rom40007FFF(addr uint16) uint8 {
	bankNumber := mbc.bank1
	if mbc.allowedRomBank2 {
		bankNumber |= mbc.bank2 << 5
	}
	return mbc.memoryBankController1RomByBankNumber(bankNumber, addr)
}

func (mbc *MBC1) memoryBankController1GetRamAddressA000BFFF(addr uint16) uint16 {
	addr &= 0x1fff
	bankAddress := addr
	if mbc.memoryModel == 1 && mbc.allowedRamBank2 {
		bankAddress |= uint16(mbc.bank2) << 13
	}
	return bankAddress
}

func (mbc *MBC1) Read(addr uint16) uint8 {
	// In its default configuration, MBC1 supports up to 512 KiB ROM with up to 32 KiB of banked RAM.

	switch addr >> 13 {
	case 0, 1: // 0x0000–0x3FFF
		return mbc.memoryBankController1RomAddress00003FFF(addr)
	case 2, 3: // 0x4000 - 0x7FFF
		return mbc.memoryBankController1Rom40007FFF(addr)
	case 5: // 0xA000 - 0xBFFF
		// This area is used to address external RAM in the cartridge (if any).
		// The RAM is only accessible if RAM is enabled,
		// otherwise reads return open bus values (often $FF, but not guaranteed) and writes are ignored.
		// Available RAM sizes are 8 KiB (at $A000–BFFF) and 32 KiB (in form of four 8K banks at $A000–BFFF).
		// 32 KiB is only available in cartridges with ROM <= 512 KiB.
		if mbc.enableRamBank {
			bankAddress := mbc.memoryBankController1GetRamAddressA000BFFF(addr)
			return mbc.ram[bankAddress]
		} else {
			return 0xFF
		}
	}

	return 0
}

func (mbc *MBC1) Write(addr uint16, val uint8) {
	// 4 "Registers"
	// - 0x0000-0x1FFF - RAMG - MBC1 RAM gate register
	// - 0x2000-0x3FFF - BANK1 - MBC1 bank register 1
	// - 0x4000-0x5FFF - BANK2 - MBC1 bank register 2
	// - 0x6000-0x7FFF - MODE - MBC1 mode register

	switch addr >> 13 {
	case 0: // 0x0000–0x1FFF, RAMG - MBC1 RAM gate register
		if val&0x0F == 0x0A {
			mbc.enableRamBank = true
		} else {
			mbc.enableRamBank = false
		}
	case 1: // 0x2000-0x3FFF, BANK1 - MBC1 bank register 1
		// If the main 5-bit ROM banking register is 0, it reads the bank as if it was set to 1.
		mbc.bank1 = val & 0x1f

		// If this register is set to $00, it behaves as if it is set to $01.
		// This means you cannot duplicate bank $00
		// into both the 0000–3FFF and 4000–7FFF ranges by setting this register to $00.
		// Even with smaller ROMs that use less than 5 bits for bank selection,
		// the full 5-bit register is still compared for the bank 00→01 translation logic.
		// As a result if the ROM is 256 KiB or smaller,
		// it is possible to map bank 0 to the 4000–7FFF region —
		// by setting the 5th bit to 1 it will prevent the 00→01 translation
		// (which looks at the full 5-bit register, and sees the value $10, not $00),
		// while the bits actually used for bank selection (4, in this example) are all 0, so bank $00 is selected.
		if mbc.bank1 == 0 {
			mbc.bank1 = 1
		}

		// If the ROM Bank Number is set to a higher value than the number of banks in the cart,
		// the bank number is masked to the required number of bits.
		// e.g. a 256 KiB cart only needs a 4-bit bank number to address all of its 16 banks,
		// so this register is masked to 4 bits. The upper bit would be ignored for bank selection.
		mbc.bank1 &= uint8(mbc.RomBanks) - 1
	case 2: // 0x4000 - 0x5FFF, BANK2 - MBC1 bank register 2
		// 1 MiB ROM or larger carts only or 32 KiB ram carts only
		if mbc.allowedRomBank2 || mbc.allowedRamBank2 {
			mbc.bank2 = val & 0x03
			// Max size number of banks
			mbc.bank2 &= (uint8(mbc.RomBanks) >> 5) - 1
		}
	case 3: // 0x6000-7FFF, MODE - MBC1 mode register
		if GetBit(val, 0) {
			mbc.memoryModel = 1
		} else {
			mbc.memoryModel = 0
		}
	case 5: // 0xA000 - 0xBFFF
		// This area is used to address external RAM in the cartridge (if any).
		// The RAM is only accessible if RAM is enabled,
		// otherwise reads return open bus values (often $FF, but not guaranteed) and writes are ignored.
		// Available RAM sizes are 8 KiB (at $A000–BFFF) and 32 KiB (in form of four 8K banks at $A000–BFFF).
		// 32 KiB is only available in cartridges with ROM <= 512 KiB.
		if mbc.enableRamBank {
			bankAddress := mbc.memoryBankController1GetRamAddressA000BFFF(addr)
			mbc.ram[bankAddress] = val
		}
	}
}
