package emulator

func (e *Emulator) memoryBankController3Write(addr uint16, val uint8) {
	switch addr >> 13 {
	case 0: // 0x0000-0x1FFF - RAM and Timer Enable
	case 1: // 0x2000 - 0x3FFF - ROM Bank Number (Write Only)
		// Ability to swap 64 different 16KiB banks of ROM
		var romBank uint32 = 1
		if val != 0 {
			romBank = uint32(val & 0x3F)
		}
		e.rom1Pointer = romBank << 14
	case 2: // 0x4000 - 0x5FFF - RAM Bank Number - or - RTC Register Select (Write Only)
		// 4 different of 8KiB banks of External Ram (for a total of 32KiB)
		if val <= 3 {
			e.extrambankPointer = uint32(val << 13)
		}
	case 3: // 0x6000 - 0x7FFF - Latch Clock Data (Write Only)
	case 5: // 0xA000 - 0xBFFF
		// A000-BFFF - RTC Register 08-0C (Read/Write)

		addr &= 0x1fff
		e.extrambank[e.extrambankPointer+uint32(addr)] = val
	}
}

func (e *Emulator) memoryBankController3Read(addr uint16) uint8 {
	switch addr >> 13 {
	case 0, 1: // 0x0000-0x3FFF
		return e.rom0[addr]
	case 2, 3: // 0x4000 - 0x5FFF
		return e.rom0[e.rom1Pointer+uint32(addr&0x3fff)]
	case 5: // 0xA000 - 0xBFFF
		addr &= 0x1fff
		return e.extrambank[e.extrambankPointer+uint32(addr)]
	}

	return 0
}

func (e *Emulator) memoryBankController3(addr uint16, val uint8, write bool) uint8 {
	// Max 2MByte ROM and/or 32KByte RAM and Timer

	if write {
		e.memoryBankController3Write(addr, val)
	} else {
		return e.memoryBankController3Read(addr)
	}

	return 0
}
