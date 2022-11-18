package emulator

import "github.com/veandco/go-sdl2/sdl"

type Memory struct {
	workRam        [0x4000]uint8
	videoRam       [0x2000]uint8
	io             [0x200]uint8
	bootRom        []byte
	bootRomEnabled bool

	lcdcControl LCDControl
	lcdStatus   LCDStatus

	rom      *Rom
	emulator *Emulator
}

func (mem *Memory) readIO(addr uint16) uint8 {
	switch addr {
	case 0xff00:
		if (^mem.io[256] & 16) != 0 {
			keyboardState := mem.emulator.window.GetKeyboardState()
			return ^(16 + keyboardState[sdl.SCANCODE_DOWN]*8 +
				keyboardState[sdl.SCANCODE_UP]*4 +
				keyboardState[sdl.SCANCODE_LEFT]*2 +
				keyboardState[sdl.SCANCODE_RIGHT])
		}
		if (^mem.io[256] & 32) != 0 {
			keyboardState := mem.emulator.window.GetKeyboardState()
			return ^(32 + keyboardState[sdl.SCANCODE_RETURN]*8 +
				keyboardState[sdl.SCANCODE_TAB]*4 +
				keyboardState[sdl.SCANCODE_Z]*2 +
				keyboardState[sdl.SCANCODE_X])
		}
		return 0xFF
	}
	ioAddr := addr & 0x1ff
	return mem.io[ioAddr]
}

func maskUnusedBits(addr uint16, val uint8) uint8 {
	switch addr {
	case 0xFF02: // SC: Serial transfer control
		val |= 0b01111110 // Unused bits returns 1s
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
	case 0xFF41: // STAT: LCD status
		val |= 0b10000000 // Unused bits returns 1s
	case 0xFF50:
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

	return val
}

func (mem *Memory) writeIO(addr uint16, val uint8) {
	val = maskUnusedBits(addr, val)

	switch addr {
	case 0xFF04: // DIV
		// When writing to DIV, if the current output is 1 and timer is enabled,
		// as the new value after reseting DIV will be 0,
		// the falling edge detector will detect a falling edge and TIMA will increase.
		if mem.emulator.timer.isFallingEdgeWritingDIV() {
			mem.emulator.timer.increaseTIMA(1, false)
		}

		// Reset Div timer
		val = 0
		mem.emulator.timer.internalTimer = 0
	case 0xFF05: // TIMA: Timer counter
		if mem.emulator.timer.timaUpdateWithTMADelayedCycles == mem.emulator.cycles {
			// If you write to TIMA during the cycle that TMA is being loaded to it [B],
			// the write will be ignored and TMA value will be written to TIMA instead.
			val = mem.GetTMA()
		} else {
			// During the strange cycle [A] you can prevent the IF flag from being set
			// and prevent the TIMA from being reloaded from TMA by writing a value to TIMA.
			// That new value will be the one that stays in the TIMA register after the instruction.
			mem.emulator.timer.timaUpdateWithTMADelayedCycles = 0
		}
	case 0xFF06: // TMA: Timer modulo
		if mem.emulator.timer.timaUpdateWithTMADelayedCycles == mem.emulator.cycles {
			mem.SetTIMA(val)
		}
	case 0xFF07: // TAC: Timer Control
		// When writing to TAC, if the previously selected multiplexer input was 1
		// and the new input is 0, TIMA will increase too.
		// This doesn't happen when the timer is disabled,
		// but it also happens when disabling the timer
		// (the same effect as writing to DIV).
		if mem.emulator.timer.isFallingEdgeWritingTAC(val) {
			// Writing to DIV, TAC or other registers won't prevent the IF flag from being set or
			// TIMA from being reloaded.
			mem.emulator.timer.increaseTIMA(1, false)
		}
	case 0xFF40: // LCDC: LCD control
		mem.SetLCDC(val)
	case 0xFF41: // STAT: LCD status
		mem.SetLCDStatus(val)
	case 0xFF46:
		for y := WIDTH - 1; y >= 0; y-- {
			mem.io[y] = mem.read8(uint16(val)<<8 | uint16(y))
		}
	case 0xFF50:
		mem.bootRomEnabled = false
	}
	ioAddr := addr & 0x1ff
	mem.io[ioAddr] = val
}

func (mem *Memory) read8(addr uint16) uint8 {
	mem.emulator.tick()

	if addr == 0xD800 {
		//val = 00
		//fmt.Println("ENTRA")
	}

	switch addr >> 13 {
	case 0, 1: // 0x0000 - 0x3FFF, ROM Bank 00
		if mem.bootRomEnabled && addr <= 0xFF {
			return mem.bootRom[addr]
		}
		value := mem.rom.controller.Read(addr)
		return mem.emulator.getGameGenieValue(addr, value)
	case 2, 3: // 0x4000 - 0x7FFF, ROM Bank 01~NN
		value := mem.rom.controller.Read(addr)
		return mem.emulator.getGameGenieValue(addr, value)
	case 4: // 0x8000 - 0x9FFF, Video RAM (VRAM)
		addr &= 0x1fff
		return mem.videoRam[addr]
	case 5: // 0xA000 - 0xBFFF, External RAM
		return mem.rom.controller.Read(addr)
	case 6: // 0xC000 - 0xDFFF, Work RAM (WRAM)
		addr &= 0x3fff
		return mem.workRam[addr]
	case 7: // 0xE000 - 0xFFFF
		if addr <= 0xFDFF { // Echo internal RAM, 0xE000-0xFDFF
			addr &= 0x3fff
			return mem.workRam[addr]
		} else { // 0xFE00-0xFFFF, OAM, I/O, High RAM (HRAM) and Interrupt Enable register (IE)
			return mem.readIO(addr)
		}
	}

	return 0
}

func (mem *Memory) write8(addr uint16, val uint8) {
	mem.emulator.tick()

	switch addr >> 13 {
	case 0, 1: // 0x0000 - 0x3FFF, ROM Bank 00
		if mem.bootRomEnabled && addr <= 0xFF {
			return
		}
		mem.rom.controller.Write(addr, val)
	case 2, 3: // 0x4000 - 0x7FFF, ROM Bank 01~NN
		mem.rom.controller.Write(addr, val)
	case 4: // 0x8000 - 0x9FFF, Video RAM (VRAM)
		addr &= 0x1fff
		mem.videoRam[addr] = val
	case 5: // 0xA000 - 0xBFFF, External RAM
		mem.rom.controller.Write(addr, val)
	case 6: // 0xC000 - 0xDFFF, Work RAM (WRAM)
		addr &= 0x3fff
		mem.workRam[addr] = val
	case 7: // 0xE000 - 0xFFFF
		if addr <= 0xFDFF { // Echo internal RAM, 0xE000-0xFDFF
			addr &= 0x3fff
			mem.workRam[addr] = val
		} else { // 0xFE00-0xFFFF, OAM, I/O, High RAM (HRAM) and Interrupt Enable register (IE)
			mem.writeIO(addr, val)
		}
	}
}

func (mem *Memory) read16(addr uint16) uint16 {
	tmp8 := mem.read8(addr)
	addr++
	result := mem.read8(addr)
	return uint16(result)<<8 | uint16(tmp8)
}

func (mem *Memory) write16(addr uint16, val uint16) {
	mem.write8(addr, uint8(val))
	addr++
	mem.write8(addr, uint8(val>>8))
}
