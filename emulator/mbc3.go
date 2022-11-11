package emulator

type ClockCounter struct {
	RTCSeconds         uint8
	RTCMinutes         uint8
	RTCHours           uint8
	RTCDL              uint8
	RTCDH              uint8
	RTCDayCounter      uint16
	Halt               bool
	DayCounterCarryBit bool
}

func (e *Emulator) memoryBankController3Write(addr uint16, val uint8) {
	switch addr >> 13 {
	case 0: // 0x0000-0x1FFF - RAM and Timer Enable
		if val == 0x0A {
			e.mbc3EnableRamBank = true
		} else {
			e.mbc3EnableRamBank = false
		}
	case 1: // 0x2000 - 0x3FFF - ROM Bank Number (Write Only)
		e.mbc3RomBank = val & 0x7F
		if e.mbc3RomBank == 0 {
			e.mbc3RomBank = 1
		}
		e.mbc3RomBank &= uint8(e.romHeader.RomBanks) - 1
	case 2: // 0x4000 - 0x5FFF - RAM Bank Number - or - RTC Register Select (Write Only)
		// 4 different of 8KiB banks of External Ram (for a total of 32KiB)
		if val <= 3 {
			// Ram Bank
			e.mbc3RamBank = val
			e.mbc3RamBank &= uint8(e.romHeader.RamBanks) - 1
			e.mbc3RegisterSelect = 0
		} else if val >= 0x08 && val <= 0xC {
			// RTC Register Activate
			e.mbc3RegisterSelect = val
		} else {
			e.mbc3RegisterSelect = 0
		}
	case 3: // 0x6000 - 0x7FFF - Latch Clock Data (Write Only)
		e.mbc3LatchRegister = (e.mbc3LatchRegister << 8) | uint16(val)
		if e.mbc3LatchRegister == 0x0001 {
			// Toggle Latched Clock
			e.mbc3LatchClock = !e.mbc3LatchClock
			e.mbc3LatchRegister = 0xFFFF
		}
	case 5: // 0xA000 - 0xBFFF
		if e.mbc3EnableRamBank {
			if e.mbc3RegisterSelect != 0 {
				switch e.mbc3RegisterSelect {
				case 0x08:
					e.mbc3ClockCounter.RTCSeconds = val
				case 0x09:
					e.mbc3ClockCounter.RTCMinutes = val
				case 0x0A:
					e.mbc3ClockCounter.RTCHours = val
				case 0x0B:
					e.mbc3ClockCounter.RTCDL = val
					e.mbc3ClockCounter.RTCDayCounter = (e.mbc3ClockCounter.RTCDayCounter & 0x0100) | uint16(val)
				case 0x0C:
					e.mbc3ClockCounter.RTCDH = val
					e.mbc3ClockCounter.RTCDayCounter = uint16(val&0x1) | (e.mbc3ClockCounter.RTCDayCounter & 0x00FF)
					e.mbc3ClockCounter.Halt = GetBit(val, 6)
					e.mbc3ClockCounter.DayCounterCarryBit = GetBit(val, 7)
				}
			} else {
				// A000-BFFF - RTC Register 08-0C (Read/Write)
				e.extrambank[(uint32(e.mbc3RamBank)<<13)+uint32(addr&0x1fff)] = val
			}
		}
	}
}

func (e *Emulator) memoryBankController3Read(addr uint16) uint8 {
	switch addr >> 13 {
	case 0, 1: // 0x0000-0x3FFF
		return e.rom0[addr]
	case 2, 3: // 0x4000 - 0x7FFF
		return e.rom0[(uint32(e.mbc3RomBank)<<14)+uint32(addr&0x3fff)]
	case 5: // 0xA000 - 0xBFFF
		if e.mbc3EnableRamBank {
			if e.mbc3RegisterSelect != 0 {
				switch e.mbc3RegisterSelect {
				case 0x08:
					return e.mbc3ClockCounter.RTCSeconds
				case 0x09:
					return e.mbc3ClockCounter.RTCMinutes
				case 0x0A:
					return e.mbc3ClockCounter.RTCHours
				case 0x0B:
					return e.mbc3ClockCounter.RTCDL
				case 0x0C:
					return e.mbc3ClockCounter.RTCDH
				}
			} else {
				return e.extrambank[(uint32(e.mbc3RamBank)<<13)+uint32(addr&0x1fff)]
			}
		}
		return 0xFF
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
