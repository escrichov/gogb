package emulator

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

func (e *Emulator) memoryBankController0(addr uint16, val uint8, write bool) uint8 {
	switch addr >> 13 {
	case 0, 1, 2, 3: // 0x2000 - 0x7FFF
		return e.rom0[addr]
	case 5: // 0xA000 - 0xBFFF
		return 0
	}

	return 0
}

func (e *Emulator) memoryBankController1(addr uint16, val uint8, write bool) uint8 {
	// In its default configuration, MBC1 supports up to 512 KiB ROM with up to 32 KiB of banked RAM.
	switch addr >> 13 {
	case 0: // 0x0000–1FFF
		if write {
			if e.mbc1MemoryModel == 1 {
				if val&0x0F == 0x0A {
					e.mbc1EnableRamBank = true
				} else {
					e.mbc1EnableRamBank = false
				}
			}
		}
		return e.rom0[addr]
	case 1: // 2000-3FFF
		if write {
			romBank := val & 0x1f
			if romBank == 0 {
				romBank = 1
			}
			e.rom1Pointer = uint32(romBank) << 14
		}
		return e.rom0[addr]
	case 2: // 0x4000 - 0x5FFF
		// This area may contain any of the further 16 KiB banks of the ROM
		// If the main 5-bit ROM banking register is 0, it reads the bank as if it was set to 1.
		if write {
			bank := val & 0x03
			if e.mbc1MemoryModel == 1 {
				e.extrambankPointer = uint32(bank)
			} else {

				e.extrambankPointer = 0
				// TODO: will set the two most significant ROM address lines.
			}
		}

		bankAddr := e.rom1Pointer + uint32(addr&0x3fff)
		return e.rom0[bankAddr]
	case 3: // 0x6000-7FFF
		if write {
			if GetBit(val, 0) {
				// S = 0 selects 16Mb/8KB mode
				e.mbc1MemoryModel = 0
			} else {
				// S = 1 selects 4Mb/32KB mode.
				e.mbc1MemoryModel = 1
			}
		}
		bankAddr := e.rom1Pointer + uint32(addr&0x3fff)
		return e.rom0[bankAddr]
	case 5: // 0xA000 - 0xBFFF
		// This area is used to address external RAM in the cartridge (if any).
		// The RAM is only accessible if RAM is enabled,
		// otherwise reads return open bus values (often $FF, but not guaranteed) and writes are ignored.
		// Available RAM sizes are 8 KiB (at $A000–BFFF) and 32 KiB (in form of four 8K banks at $A000–BFFF).
		// 32 KiB is only available in cartridges with ROM <= 512 KiB.
		if e.mbc1EnableRamBank {
			addr &= 0x1fff
			if write {
				e.extrambank[e.extrambankPointer+uint32(addr)] = val
			}
			return e.extrambank[e.extrambankPointer+uint32(addr)]
		}
		return 0
	}

	return 0
}

func (e *Emulator) memoryBankController2(addr uint16, val uint8, write bool) uint8 {
	return 0
}

func (e *Emulator) memoryBankController3(addr uint16, val uint8, write bool) uint8 {
	switch addr >> 13 {
	case 0: // 0x0000 - 0x1FFF
		return e.rom0[addr]
	case 1: // 0x2000 - 0x3FFF
		if write {
			// Ability to swap 64 different 16KiB banks of ROM
			var romBank uint32 = 1
			if val != 0 {
				romBank = uint32(val & 0x3F)
			}
			e.rom1Pointer = romBank << 14
		}

		return e.rom0[addr]
	case 2: // 0x4000 - 0x5FFF
		// 4 different of 8KiB banks of External Ram (for a total of 32KiB)
		if write && val <= 3 {
			e.extrambankPointer = uint32(val << 13)
		}
		return e.rom0[e.rom1Pointer+uint32(addr&0x3fff)]
	case 3: // 0x6000 - 0x7FFF
		return e.rom0[e.rom1Pointer+uint32(addr&0x3fff)]
	case 5: // 0xA000 - 0xBFFF
		addr &= 0x1fff
		if write {
			e.extrambank[e.extrambankPointer+uint32(addr)] = val
		}
		return e.extrambank[e.extrambankPointer+uint32(addr)]
	}

	return 0
}

func (e *Emulator) memoryBankController5(addr uint16, val uint8, write bool) uint8 {
	switch addr >> 13 {
	case 0, 2, 3, 5:
		return e.memoryBankController3(addr, val, write)
	case 1: // 0x2000 - 0x3FFF
		if write {
			// Ability to swap 64 different 16KiB banks of ROM
			if addr <= 0x2FFF {
				var romBank = uint32(val & 0x3F)
				e.rom1Pointer = romBank << 14
			} else {
				// TODO: Implement set bit 9
			}
		}
		return e.rom0[addr]
	}

	return 0
}

func (e *Emulator) memoryBankController6(addr uint16, val uint8, write bool) uint8 {
	return 0
}

func (e *Emulator) memoryBankController7(addr uint16, val uint8, write bool) uint8 {
	return 0
}

func (e *Emulator) memMemoryBankController(addr uint16, val uint8, write bool) uint8 {
	switch e.memoryBankController {
	case 0:
		return e.memoryBankController0(addr, val, write)
	case 1:
		return e.memoryBankController1(addr, val, write)
	case 3:
		return e.memoryBankController3(addr, val, write)
	case 5:
		return e.memoryBankController3(addr, val, write)
	default:
		log.Fatal("Unsupported memory bank controller: ", e.memoryBankController)
	}

	return 0
}

func (e *Emulator) mem8(addr uint16, val uint8, write bool) uint8 {
	e.tick()

	if addr == 0xD800 {
		val = 00
		fmt.Println("ENTRA")
	}

	switch addr >> 13 {
	case 0: // 0x0000 - 0x1FFF
		if e.bootRomEnabled && addr <= 0xFF {
			return e.bootRom[addr]
		}
		return e.memMemoryBankController(addr, val, write)
	case 1, 2, 3, 5: // 0x2000 - 0xBFFF
		return e.memMemoryBankController(addr, val, write)
	case 4: // 0x8000 - 0x9FFF
		addr &= 8191
		if write {
			e.videoRam[addr] = val
		}
		return e.videoRam[addr]
	case 7: // 0xE000 - 0xFFFF
		if addr >= 0xFE00 {
			if write {
				switch addr {
				case 0xFF02: // SC: Serial transfer control
					val |= 0b01111110 // Unused bits returns 1s
				case 0xFF04: // DIV
					val = 0
					e.divTimer = 0
				case 0xFF05: // TIMA
					e.timaTimer = int(val)
				case 0xFF07: // TAC: Timer Control
					val |= 0b11111000 // Unused bits returns 1s
				case 0xFF0F: // IF
					val |= 0b11100000 // Unused bits returns 1s
				case 0xFF10: // Audio NR10: Channel 1 sweep
					val |= 0b10000000 // Unused bits returns 1s
				case 0xFF1C: // NR32: Channel 3 output level
					val |= 0b10011111 // Unused bits returns 1s
				case 0xFF1A: // Audio NR30: Channel 3 DAC enable
					val |= 0b01111111 // Unused bits returns 1s
				case 0xFF20: // NR41: Channel 4 length timer [write-only]
					val |= 0b11000000 // Unused bits returns 1s
				case 0xFF23: // NR44: Channel 4 control
					val |= 0b00111111 // Unused bits returns 1s
				case 0xFF26: // NR52: Sound on/off
					val |= 0b01110000 // Unused bits returns 1s
				case 0xFF40: // LCDC: LCD control
					e.SetLCDC(val)
				case 0xFF41: // STAT: LCD status
					val |= 0b10000000 // Unused bits returns 1s
				case 0xFF46:
					for y := WIDTH - 1; y >= 0; y-- {
						e.io[y] = e.read8(uint16(val)<<8 | uint16(y))
					}
				case 0xFF50:
					e.bootRomEnabled = false
					val |= 0b11111111 // Unused bits returns 1s
				case 0xFF03, 0xFF08, 0xFF09, 0xFF0A, 0xFF0B,
					0xFF0C, 0xFF0D, 0xFF0E, 0xFF15, 0xFF1F,
					0xFF27, 0xFF28, 0xFF29,
					0xFF4C, 0xFF4D, 0xFF4E, 0xFF4F,
					0xFF51, 0xFF52, 0xFF53, 0xFF54,
					0xFF55, 0xFF56, 0xFF57, 0xFF58,
					0xFF59, 0xFF5A, 0xFF5B, 0xFF5C,
					0xFF5D, 0xFF5E, 0xFF5F,
					0xFF60, 0xFF61, 0xFF62, 0xFF63, 0xFF64,
					0xFF65, 0xFF66, 0xFF67, 0xFF68,
					0xFF69, 0xFF6A, 0xFF6B, 0xFF6C,
					0xFF6D, 0xFF6E, 0xFF6F,
					0xFF70, 0xFF71, 0xFF72, 0xFF73, 0xFF74,
					0xFF75, 0xFF76, 0xFF77, 0xFF78,
					0xFF79, 0xFF7A, 0xFF7B, 0xFF7C,
					0xFF7D, 0xFF7E, 0xFF7F: // Unused
					val |= 0b11111111 // Unused bits returns 1s
				}
				ioAddr := addr & 0x1ff
				e.io[ioAddr] = val
			}

			switch addr {
			case 0xff00:
				if (^e.io[256] & 16) != 0 {
					return ^(16 + e.keyboardState[sdl.SCANCODE_DOWN]*8 +
						e.keyboardState[sdl.SCANCODE_UP]*4 +
						e.keyboardState[sdl.SCANCODE_LEFT]*2 +
						e.keyboardState[sdl.SCANCODE_RIGHT])
				}
				if (^e.io[256] & 32) != 0 {
					return ^(32 + e.keyboardState[sdl.SCANCODE_RETURN]*8 +
						e.keyboardState[sdl.SCANCODE_TAB]*4 +
						e.keyboardState[sdl.SCANCODE_Z]*2 +
						e.keyboardState[sdl.SCANCODE_X])
				}
				return 0xFF
			}
			ioAddr := addr & 0x1ff
			return e.io[ioAddr]
		} else { // Echo internal RAM
			addr &= 0x3fff
			if write {
				e.workRam[addr] = val
			}
			return e.workRam[addr]
		}
	case 6: // 0xC000 - 0xDFFF, Internal RAM
		addr &= 0x3fff
		if write {
			e.workRam[addr] = val
		}
		return e.workRam[addr]
	}

	return 0
}

func (e *Emulator) read16(addr uint16) uint16 {
	tmp8 := e.mem8(addr, 0, false)
	addr++
	result := e.mem8(addr, 0, false)
	return uint16(result)<<8 | uint16(tmp8)
}

func (e *Emulator) read8(addr uint16) uint8 {
	return e.mem8(addr, 0, false)
}

func (e *Emulator) write16(addr uint16, val uint16) {
	e.mem8(addr, uint8(val), true)
	addr++
	e.mem8(addr, uint8(val>>8), true)
}

func (e *Emulator) write8(addr uint16, val uint8) {
	e.mem8(addr, val, true)
}
