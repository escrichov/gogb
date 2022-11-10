package emulator

func (e *Emulator) memoryBankController1RomByBankNumber(bankNumber uint8, addr uint16) uint8 {
	bankAddr := uint32(bankNumber)<<14 + uint32(addr&0x3fff)
	return e.rom0[bankAddr]
}

func (e *Emulator) memoryBankController1RomAddress00003FFF(addr uint16) uint8 {
	bankNumber := uint8(0)
	if e.mbc1MemoryModel == 1 {
		if e.mbc1AllowedRomBank2 {
			bankNumber = e.mbc1Bank2 << 5
		}
	}
	return e.memoryBankController1RomByBankNumber(bankNumber, addr)
}

func (e *Emulator) memoryBankController1Rom40007FFF(addr uint16) uint8 {
	bankNumber := e.mbc1Bank1
	if e.mbc1AllowedRomBank2 {
		bankNumber |= e.mbc1Bank2 << 5
	}
	return e.memoryBankController1RomByBankNumber(bankNumber, addr)
}

func (e *Emulator) memoryBankController1GetRamAddressA000BFFF(addr uint16) uint16 {
	addr &= 0x1fff
	bankAddress := addr
	if e.mbc1MemoryModel == 1 && e.mbc1AllowedRamBank2 {
		bankAddress |= uint16(e.mbc1Bank2) << 13
	}
	return bankAddress
}

func (e *Emulator) memoryBankController1Write(addr uint16, val uint8) {
	// 4 "Registers"
	// - 0x0000-0x1FFF - RAMG - MBC1 RAM gate register
	// - 0x2000-0x3FFF - BANK1 - MBC1 bank register 1
	// - 0x4000-0x5FFF - BANK2 - MBC1 bank register 2
	// - 0x6000-0x7FFF - MODE - MBC1 mode register

	switch addr >> 13 {
	case 0: // 0x0000–0x1FFF, RAMG - MBC1 RAM gate register
		if val&0x0F == 0x0A {
			e.mbc1EnableRamBank = true
		} else {
			e.mbc1EnableRamBank = false
		}
	case 1: // 0x2000-0x3FFF, BANK1 - MBC1 bank register 1
		// If the main 5-bit ROM banking register is 0, it reads the bank as if it was set to 1.
		e.mbc1Bank1 = val & 0x1f

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
		if e.mbc1Bank1 == 0 {
			e.mbc1Bank1 = 1
		}

		// If the ROM Bank Number is set to a higher value than the number of banks in the cart,
		// the bank number is masked to the required number of bits.
		// e.g. a 256 KiB cart only needs a 4-bit bank number to address all of its 16 banks,
		// so this register is masked to 4 bits. The upper bit would be ignored for bank selection.
		e.mbc1Bank1 &= uint8(e.romHeader.RomBanks) - 1
	case 2: // 0x4000 - 0x5FFF, BANK2 - MBC1 bank register 2
		// 1 MiB ROM or larger carts only or 32 KiB ram carts only
		if e.mbc1AllowedRomBank2 || e.mbc1AllowedRamBank2 {
			e.mbc1Bank2 = val & 0x03
			// Max size number of banks
			e.mbc1Bank2 &= (uint8(e.romHeader.RomBanks) >> 5) - 1
		}
	case 3: // 0x6000-7FFF, MODE - MBC1 mode register
		if GetBit(val, 0) {
			e.mbc1MemoryModel = 1
		} else {
			e.mbc1MemoryModel = 0
		}
	case 5: // 0xA000 - 0xBFFF
		// This area is used to address external RAM in the cartridge (if any).
		// The RAM is only accessible if RAM is enabled,
		// otherwise reads return open bus values (often $FF, but not guaranteed) and writes are ignored.
		// Available RAM sizes are 8 KiB (at $A000–BFFF) and 32 KiB (in form of four 8K banks at $A000–BFFF).
		// 32 KiB is only available in cartridges with ROM <= 512 KiB.
		if e.mbc1EnableRamBank {
			bankAddress := e.memoryBankController1GetRamAddressA000BFFF(addr)
			e.extrambank[bankAddress] = val
		}
	}
}

func (e *Emulator) memoryBankController1Read(addr uint16) uint8 {
	// In its default configuration, MBC1 supports up to 512 KiB ROM with up to 32 KiB of banked RAM.

	switch addr >> 13 {
	case 0, 1: // 0x0000–0x3FFF
		return e.memoryBankController1RomAddress00003FFF(addr)
	case 2, 3: // 0x4000 - 0x7FFF
		return e.memoryBankController1Rom40007FFF(addr)
	case 5: // 0xA000 - 0xBFFF
		// This area is used to address external RAM in the cartridge (if any).
		// The RAM is only accessible if RAM is enabled,
		// otherwise reads return open bus values (often $FF, but not guaranteed) and writes are ignored.
		// Available RAM sizes are 8 KiB (at $A000–BFFF) and 32 KiB (in form of four 8K banks at $A000–BFFF).
		// 32 KiB is only available in cartridges with ROM <= 512 KiB.
		if e.mbc1EnableRamBank {
			bankAddress := e.memoryBankController1GetRamAddressA000BFFF(addr)
			return e.extrambank[bankAddress]
		} else {
			return 0xFF
		}
	}

	return 0
}

func (e *Emulator) memoryBankController1(addr uint16, val uint8, write bool) uint8 {
	// In its default configuration, MBC1 supports up to 512 KiB ROM with up to 32 KiB of banked RAM.
	// 4 "Registers"
	// - 0x0000-0x1FFF - RAMG - MBC1 RAM gate register
	// - 0x2000-0x3FFF - BANK1 - MBC1 bank register 1
	// - 0x4000-0x5FFF - BANK2 - MBC1 bank register 2
	// - 0x6000-0x7FFF - MODE - MBC1 mode register

	if write {
		e.memoryBankController1Write(addr, val)
	} else {
		return e.memoryBankController1Read(addr)
	}

	return 0
}
