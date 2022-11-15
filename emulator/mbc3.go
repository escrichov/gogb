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

type MBC3 struct {
	*BaseMBC

	enableRamBank  bool
	enableRTC      bool
	romBank        uint8
	ramBank        uint8
	clockCounter   ClockCounter
	latchRegister  uint16
	latchClock     bool
	registerSelect uint8
}

func NewMBC3(baseMBC *BaseMBC) *MBC3 {
	mbc := &MBC3{
		BaseMBC:       baseMBC,
		romBank:       1,
		latchRegister: 0xFFFF,
	}

	return mbc
}

func (mbc *MBC3) Write(addr uint16, val uint8) {
	switch addr >> 13 {
	case 0: // 0x0000-0x1FFF - RAM and Timer Enable
		if val == 0x0A {
			mbc.enableRamBank = true
		} else {
			mbc.enableRamBank = false
		}
	case 1: // 0x2000 - 0x3FFF - ROM Bank Number (Write Only)
		mbc.romBank = val & 0x7F
		if mbc.romBank == 0 {
			mbc.romBank = 1
		}
		mbc.romBank &= uint8(mbc.RomBanks) - 1
	case 2: // 0x4000 - 0x5FFF - RAM Bank Number - or - RTC Register Select (Write Only)
		// 4 different of 8KiB banks of External Ram (for a total of 32KiB)
		if val <= 3 {
			// Ram Bank
			mbc.ramBank = val
			mbc.ramBank &= uint8(mbc.RamBanks) - 1
			mbc.registerSelect = 0
		} else if val >= 0x08 && val <= 0xC {
			// RTC Register Activate
			mbc.registerSelect = val
		} else {
			mbc.registerSelect = 0
		}
	case 3: // 0x6000 - 0x7FFF - Latch Clock Data (Write Only)
		mbc.latchRegister = (mbc.latchRegister << 8) | uint16(val)
		if mbc.latchRegister == 0x0001 {
			// Toggle Latched Clock
			mbc.latchClock = !mbc.latchClock
			mbc.latchRegister = 0xFFFF
		}
	case 5: // 0xA000 - 0xBFFF
		if mbc.enableRamBank {
			if mbc.registerSelect != 0 {
				switch mbc.registerSelect {
				case 0x08:
					mbc.clockCounter.RTCSeconds = val
				case 0x09:
					mbc.clockCounter.RTCMinutes = val
				case 0x0A:
					mbc.clockCounter.RTCHours = val
				case 0x0B:
					mbc.clockCounter.RTCDL = val
					mbc.clockCounter.RTCDayCounter = (mbc.clockCounter.RTCDayCounter & 0x0100) | uint16(val)
				case 0x0C:
					mbc.clockCounter.RTCDH = val
					mbc.clockCounter.RTCDayCounter = uint16(val&0x1) | (mbc.clockCounter.RTCDayCounter & 0x00FF)
					mbc.clockCounter.Halt = GetBit(val, 6)
					mbc.clockCounter.DayCounterCarryBit = GetBit(val, 7)
				}
			} else {
				// A000-BFFF - RTC Register 08-0C (Read/Write)
				mbc.ram[(uint32(mbc.ramBank)<<13)+uint32(addr&0x1fff)] = val
			}
		}
	}
}

func (mbc *MBC3) Read(addr uint16) uint8 {
	switch addr >> 13 {
	case 0, 1: // 0x0000-0x3FFF
		return mbc.rom[addr]
	case 2, 3: // 0x4000 - 0x7FFF
		return mbc.rom[(uint32(mbc.romBank)<<14)+uint32(addr&0x3fff)]
	case 5: // 0xA000 - 0xBFFF
		if mbc.enableRamBank {
			if mbc.registerSelect != 0 {
				switch mbc.registerSelect {
				case 0x08:
					return mbc.clockCounter.RTCSeconds
				case 0x09:
					return mbc.clockCounter.RTCMinutes
				case 0x0A:
					return mbc.clockCounter.RTCHours
				case 0x0B:
					return mbc.clockCounter.RTCDL
				case 0x0C:
					return mbc.clockCounter.RTCDH
				}
			} else {
				return mbc.ram[(uint32(mbc.ramBank)<<13)+uint32(addr&0x1fff)]
			}
		}
		return 0xFF
	}

	return 0
}
