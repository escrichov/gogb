package emulator

import (
	"log"
)

func (e *Emulator) instRet() {
	addr := e.pop()
	e.cpu.PC = addr
	e.tick()

	// Return to halt after interrupt
	// The behavior is different when ei (whose effect is typically delayed by one instruction)
	// is followed immediately by a halt, and an interrupt is pending as the halt is executed.
	// The interrupt is serviced and the handler called,
	// but the interrupt returns to the halt, which is executed again,
	// and thus waits for another interrupt.
	if e.isHaltBugEIActive {
		e.halt = 1
		e.isHaltBugEIActive = false
	}
}

func (e *Emulator) instCall(addr uint16) {
	e.push(e.cpu.PC)
	e.cpu.PC = addr
}

func (e *Emulator) instHalt() {
	if e.mem.hasPendingInterrupts() {
		e.isInterruptPendingInFirstHaltExecution = true
	}
	if e.delayedActivateIMEatInstruction == e.numInstructions {
		e.isIMEDelayedInFirstHaltExecution = true
	}

	e.halt = 1
}

func (e *Emulator) instEI() {
	e.delayedActivateIMEatInstruction = e.numInstructions + 1
}

func (e *Emulator) CPURun() {
	opcode := e.popPC()
	e.numInstructions++

	switch opcode {
	case 0: // NOP
	case 8: // LD (u16), SP
		e.mem.write16(e.popPC16(), e.cpu.GetSP())
	case 16: // STOP
		// TODO: improve stop
		// Timing is 1 Cycle
		e.instHalt()
		e.popPC()
	case 24: // JR (unconditional)
		i8 := int8(e.popPC())
		addr := int32(e.cpu.PC) + int32(i8)
		e.cpu.PC = uint16(addr)
		e.tick()
	case 32, 40, 48, 56: // JR (conditional)
		i8 := int8(e.popPC())
		if e.cpu.checkCondition((opcode >> 3) & 0x3) {
			addr := int32(e.cpu.PC) + int32(i8)
			e.cpu.PC = uint16(addr)
			e.tick()
		}
	case 1, 17, 33, 49: // LD r16, u16
		u16 := e.popPC16()
		number := (opcode >> 4) & 0x3
		e.cpu.r16Group1Set(number, u16)
	case 9, 25, 41, 57: // ADD HL, r16
		number := (opcode >> 4) & 0x3
		r16 := e.cpu.r16Group1Get(number)
		hl := e.cpu.GetHL()
		total := int32(hl) + int32(r16)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(int32(hl&0xFFF) > (total & 0xFFF))
		e.cpu.SetCarryFlag(total > 0xFFFF)
		e.cpu.SetHL(uint16(total))
		e.tick()
	case 2, 18, 34, 50: // LD (r16), A
		number := (opcode >> 4) & 0x3
		e.mem.write8(e.cpu.r16Group2Get(number), e.cpu.GetA())
	case 10, 26, 42, 58: // LD A, (r16)
		number := (opcode >> 4) & 0x3
		e.cpu.SetA(e.mem.read8(e.cpu.r16Group2Get(number)))
	case 3, 19, 35, 51: // INC r16
		number := (opcode >> 4) & 0x3
		r16 := e.cpu.r16Group1Get(number)
		e.cpu.r16Group1Set(number, r16+1)
		e.tick()
	case 11, 27, 43, 59: // DEC r16
		number := (opcode >> 4) & 0x3
		r16 := e.cpu.r16Group1Get(number)
		e.cpu.r16Group1Set(number, r16-1)
		e.tick()
	case 4, 12, 20, 28, 36, 44, 52, 60: // INC r8
		number := (opcode >> 3) & 0x7
		r8 := e.r8Get(number)
		result := r8 + 1
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(result&0xF == 0)
		e.r8Set(number, result)
	case 5, 13, 21, 29, 37, 45, 53, 61: // DEC r8
		number := (opcode >> 3) & 0x7
		r8 := e.r8Get(number)
		result := r8 - 1
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(true)
		e.cpu.SetHalfCarryFlag(r8&0x0F == 0)
		e.r8Set(number, r8-1)
	case 6, 14, 22, 30, 38, 46, 54, 62: // LD r8, u8
		u8 := e.popPC()
		number := (opcode >> 3) & 0x7
		e.r8Set(number, u8)
	case 7: // RLCA
		a := e.cpu.GetA()
		bit7 := GetBit(a, 7)
		result := (a << 1) | BoolToUint8(bit7)
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(false)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(false)
		e.cpu.SetCarryFlag(bit7)
	case 15: // RRCA
		a := e.cpu.GetA()
		bit0 := GetBit(a, 0)
		result := (a >> 1) | (BoolToUint8(bit0) << 7)
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(false)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(false)
		e.cpu.SetCarryFlag(bit0)
	case 23: // RLA
		a := e.cpu.GetA()
		bit7 := GetBit(a, 7)
		carry := BoolToUint8(e.cpu.GetCarryFlag())
		result := (a << 1) | carry
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(false)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(false)
		e.cpu.SetCarryFlag(bit7)
	case 31: // RRA
		a := e.cpu.GetA()
		bit0 := GetBit(a, 0)
		carry := BoolToUint8(e.cpu.GetCarryFlag())
		result := (a >> 1) | (carry << 7)
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(false)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(false)
		e.cpu.SetCarryFlag(bit0)
	case 39: // DAA
		a := e.cpu.GetA()
		result := a
		if e.cpu.GetSubtractFlag() { // after a subtraction, only adjust if (half-)carry occurred
			if e.cpu.GetCarryFlag() {
				result -= 0x60
			}
			if e.cpu.GetHalfCarryFlag() {
				result -= 0x6
			}
		} else { // after an addition, adjust if (half-)carry occurred or if result is out of bounds
			if e.cpu.GetCarryFlag() || a > 0x99 {
				result += 0x60
				e.cpu.SetCarryFlag(true)
			}
			if e.cpu.GetHalfCarryFlag() || (a&0x0f) > 9 {
				result += 0x6
			}
		}
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetHalfCarryFlag(false)
	case 47: // CPL
		a := e.cpu.GetA()
		result := a ^ 0xFF
		e.cpu.SetA(result)
		e.cpu.SetSubtractFlag(true)
		e.cpu.SetHalfCarryFlag(true)
	case 55: // SCF
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(false)
		e.cpu.SetCarryFlag(true)
	case 63: // CCF
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(false)
		e.cpu.SetCarryFlag(!e.cpu.GetCarryFlag())
	case 118: // HALT
		e.instHalt()
	case 64, 65, 66, 67, 68, 69, 70, 71, // LD r8, r8
		72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87,
		88, 89, 90, 91, 92, 93, 94, 95,
		96, 97, 98, 99, 100, 101, 102, 103,
		104, 105, 106, 107, 108, 109, 110, 111,
		112, 113, 114, 115, 116, 117, 119, // 118 is missing because is defined HALT,
		120, 121, 122, 123, 124, 125, 126, 127:
		numberSource := opcode & 0x7
		numberDestination := (opcode >> 3) & 0x7
		e.r8Set(numberDestination, e.r8Get(numberSource))
	case 128, 129, 130, 131, 132, 133, 134, 135: // ADD A, r8
		number := opcode & 0x7
		r8 := e.r8Get(number)
		a := e.cpu.GetA()
		result := a + r8
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(a%16+r8%16 > 15)
		e.cpu.SetCarryFlag(uint16(a)+uint16(r8) > 255)
	case 136, 137, 138, 139, 140, 141, 142, 143: // ADC A, r8
		number := opcode & 0x7
		r8 := e.r8Get(number)
		carry := BoolToUint8(e.cpu.GetCarryFlag())
		a := e.cpu.GetA()
		result := a + r8 + carry
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(a%16+r8%16+carry > 15)
		e.cpu.SetCarryFlag(uint16(a)+uint16(r8)+uint16(carry) > 255)
	case 144, 145, 146, 147, 148, 149, 150, 151: // SUB A, r8
		number := opcode & 0x7
		r8 := e.r8Get(number)
		a := e.cpu.GetA()
		result := a - r8
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(true)
		e.cpu.SetHalfCarryFlag(a%16-r8%16 > 15)
		e.cpu.SetCarryFlag(uint16(a)-uint16(r8) > 255)
	case 152, 153, 154, 155, 156, 157, 158, 159: // SBC A, r8
		number := opcode & 0x7
		r8 := e.r8Get(number)
		carry := BoolToUint8(e.cpu.GetCarryFlag())
		a := e.cpu.GetA()
		result := a - r8 - carry
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(true)
		e.cpu.SetHalfCarryFlag(a%16-r8%16-carry > 15)
		e.cpu.SetCarryFlag(uint16(a)-uint16(r8)-uint16(carry) > 255)
	case 160, 161, 162, 163, 164, 165, 166, 167: // AND A, r8
		number := opcode & 0x7
		r8 := e.r8Get(number)
		a := e.cpu.GetA()
		result := a & r8
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(true)
		e.cpu.SetCarryFlag(false)
	case 168, 169, 170, 171, 172, 173, 174, 175: // XOR A, r8
		number := opcode & 0x7
		r8 := e.r8Get(number)
		a := e.cpu.GetA()
		result := a ^ r8
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(false)
		e.cpu.SetCarryFlag(false)
	case 176, 177, 178, 179, 180, 181, 182, 183: // OR A, r8
		number := opcode & 0x7
		r8 := e.r8Get(number)
		a := e.cpu.GetA()
		result := a | r8
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(false)
		e.cpu.SetCarryFlag(false)
	case 184, 185, 186, 187, 188, 189, 190, 191: // CP A, r8
		number := opcode & 0x7
		r8 := e.r8Get(number)
		a := e.cpu.GetA()
		result := a - r8
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(true)
		e.cpu.SetHalfCarryFlag(a%16-r8%16 > 15)
		e.cpu.SetCarryFlag(uint16(a)-uint16(r8) > 255)
	case 192, 200, 208, 216: // RET condition
		e.tick()
		if e.cpu.checkCondition((opcode >> 3) & 0x3) {
			e.instRet()
		}
	case 224: // LD (FF00 + u8), A
		u8 := e.popPC()
		addr := 0xFF00 + uint16(u8)
		e.mem.write8(addr, e.cpu.GetA())
	case 232: // ADD SP, i8
		i8 := int8(e.popPC())
		sp := e.cpu.GetSP()
		result := uint16(int32(sp) + int32(i8))
		tmpVal := sp ^ uint16(i8) ^ result
		e.cpu.SetSP(result)
		e.cpu.SetZeroFlag(false)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag((tmpVal & 0x10) == 0x10)
		e.cpu.SetCarryFlag((tmpVal & 0x100) == 0x100)
		e.tick()
		e.tick()
	case 240: // LD A, (FF00 + u8)
		u8 := e.popPC()
		addr := 0xFF00 + uint16(u8)
		e.cpu.SetA(e.mem.read8(addr))
	case 248: // LD HL, SP + i8
		i8 := int8(e.popPC())
		sp := e.cpu.GetSP()
		result := uint16(int32(sp) + int32(i8))
		tmpVal := sp ^ uint16(i8) ^ result
		e.cpu.SetHL(result)
		e.cpu.SetZeroFlag(false)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag((tmpVal & 0x10) == 0x10)
		e.cpu.SetCarryFlag((tmpVal & 0x100) == 0x100)
		e.tick()
	case 193, 209, 225, 241: // POP r16
		number := (opcode >> 4) & 0x3
		value := e.pop()
		e.cpu.r16Group3Set(number, value)
	case 201: // RET
		e.instRet()
	case 217: // RETI
		e.instEI()
		e.instRet()
	case 233: // JP HL
		e.cpu.PC = e.cpu.GetHL()
	case 249: // LD SP, HL
		e.cpu.SetSP(e.cpu.GetHL())
		e.tick()
	case 194, 202, 210, 218: // JP condition
		addr := e.popPC16()
		if e.cpu.checkCondition((opcode >> 3) & 0x3) {
			e.cpu.PC = addr
			e.tick()
		}
	case 226: // LD (FF00+C), A
		addr := 0xFF00 + uint16(e.cpu.GetC())
		e.mem.write8(addr, e.cpu.GetA())
	case 234: // LD (u16), A
		u16 := e.popPC16()
		e.mem.write8(u16, e.cpu.GetA())
	case 242: // LD A, (0xFF00+C)
		addr := 0xFF00 + uint16(e.cpu.GetC())
		e.cpu.SetA(e.mem.read8(addr))
	case 250: // LD A, (u16)
		u16 := e.popPC16()
		val := e.mem.read8(u16)
		e.cpu.SetA(val)
	case 195: // JP u16
		e.cpu.PC = e.popPC16()
		e.tick()
	case 243: // DI
		e.IME = 0
	case 251: // EI
		e.instEI()
	case 196, 204, 212, 220: // CALL condition
		u16 := e.popPC16()
		if e.cpu.checkCondition((opcode >> 3) & 0x3) {
			e.instCall(u16)
		}
	case 197, 213, 229, 245: // PUSH r16
		number := (opcode >> 4) & 0x3
		value := e.cpu.r16Group3Get(number)
		e.push(value)
	case 205: // CALL u16
		u16 := e.popPC16()
		e.instCall(u16)
	case 198: // ADD A, u8
		u8 := e.popPC()
		a := e.cpu.GetA()
		result := a + u8
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(a%16+u8%16 > 15)
		e.cpu.SetCarryFlag(uint16(a)+uint16(u8) > 255)
	case 206: // ADC A, u8
		u8 := e.popPC()
		carry := BoolToUint8(e.cpu.GetCarryFlag())
		a := e.cpu.GetA()
		result := a + u8 + carry
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(a%16+u8%16+carry > 15)
		e.cpu.SetCarryFlag(uint16(a)+uint16(u8)+uint16(carry) > 255)
	case 214: // SUB A, u8
		u8 := e.popPC()
		a := e.cpu.GetA()
		result := a - u8
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(true)
		e.cpu.SetHalfCarryFlag(a%16-u8%16 > 15)
		e.cpu.SetCarryFlag(uint16(a)-uint16(u8) > 255)
	case 222: // SBC A, u8
		u8 := e.popPC()
		carry := BoolToUint8(e.cpu.GetCarryFlag())
		a := e.cpu.GetA()
		result := a - u8 - carry
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(true)
		e.cpu.SetHalfCarryFlag(a%16-u8%16-carry > 15)
		e.cpu.SetCarryFlag(uint16(a)-uint16(u8)-uint16(carry) > 255)
	case 230: // AND A, u8
		u8 := e.popPC()
		a := e.cpu.GetA()
		result := a & u8
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(true)
		e.cpu.SetCarryFlag(false)
	case 238: // XOR A, u8
		u8 := e.popPC()
		a := e.cpu.GetA()
		result := a ^ u8
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(false)
		e.cpu.SetCarryFlag(false)
	case 246: // OR A, u8
		u8 := e.popPC()
		a := e.cpu.GetA()
		result := a | u8
		e.cpu.SetA(result)
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(false)
		e.cpu.SetHalfCarryFlag(false)
		e.cpu.SetCarryFlag(false)
	case 254: // CP A, u8
		u8 := e.popPC()
		a := e.cpu.GetA()
		result := a - u8
		e.cpu.SetZeroFlag(result == 0)
		e.cpu.SetSubtractFlag(true)
		e.cpu.SetHalfCarryFlag(a%16-u8%16 > 15)
		e.cpu.SetCarryFlag(uint16(a)-uint16(u8) > 255)
	case 199, 207, 215, 223, 231, 239, 247, 255: // RST (Call to 00EXP000)
		addr := uint16(opcode & 0x38)
		e.instCall(addr)
	case 0xCB:
		opcode := e.popPC()
		switch opcode {
		case 0, 1, 2, 3, 4, 5, 6, 7: // RLC
			number := opcode & 0x7
			r8 := e.r8Get(number)
			bit7 := GetBit(r8, 7)
			result := (r8 << 1) | BoolToUint8(bit7)
			e.r8Set(number, result)
			e.cpu.SetZeroFlag(result == 0)
			e.cpu.SetSubtractFlag(false)
			e.cpu.SetHalfCarryFlag(false)
			e.cpu.SetCarryFlag(bit7)
		case 8, 9, 10, 11, 12, 13, 14, 15: // RRC
			number := opcode & 0x7
			r8 := e.r8Get(number)
			bit0 := GetBit(r8, 0)
			result := (r8 >> 1) | (BoolToUint8(bit0) << 7)
			e.r8Set(number, result)
			e.cpu.SetZeroFlag(result == 0)
			e.cpu.SetSubtractFlag(false)
			e.cpu.SetHalfCarryFlag(false)
			e.cpu.SetCarryFlag(bit0)
		case 16, 17, 18, 19, 20, 21, 22, 23: // RL
			number := opcode & 0x7
			r8 := e.r8Get(number)
			bit7 := GetBit(r8, 7)
			carry := BoolToUint8(e.cpu.GetCarryFlag())
			result := (r8 << 1) | carry
			e.r8Set(number, result)
			e.cpu.SetZeroFlag(result == 0)
			e.cpu.SetSubtractFlag(false)
			e.cpu.SetHalfCarryFlag(false)
			e.cpu.SetCarryFlag(bit7)
		case 24, 25, 26, 27, 28, 29, 30, 31: // RR
			number := opcode & 0x7
			r8 := e.r8Get(number)
			bit0 := GetBit(r8, 0)
			carry := BoolToUint8(e.cpu.GetCarryFlag())
			result := (r8 >> 1) | (carry << 7)
			e.r8Set(number, result)
			e.cpu.SetZeroFlag(result == 0)
			e.cpu.SetSubtractFlag(false)
			e.cpu.SetHalfCarryFlag(false)
			e.cpu.SetCarryFlag(bit0)
		case 32, 33, 34, 35, 36, 37, 38, 39: // SLA
			number := opcode & 0x7
			r8 := e.r8Get(number)
			bit7 := GetBit(r8, 7)
			result := r8 << 1
			e.r8Set(number, result)
			e.cpu.SetZeroFlag(result == 0)
			e.cpu.SetSubtractFlag(false)
			e.cpu.SetHalfCarryFlag(false)
			e.cpu.SetCarryFlag(bit7)
		case 40, 41, 42, 43, 44, 45, 46, 47: // SRA
			number := opcode & 0x7
			r8 := e.r8Get(number)
			bit0 := GetBit(r8, 0)
			bit7 := GetBit(r8, 7)
			result := r8>>1 | (BoolToUint8(bit7) << 7)
			e.r8Set(number, result)
			e.cpu.SetZeroFlag(result == 0)
			e.cpu.SetSubtractFlag(false)
			e.cpu.SetHalfCarryFlag(false)
			e.cpu.SetCarryFlag(bit0)
		case 48, 49, 50, 51, 52, 53, 54, 55: // SWAP
			number := opcode & 0x7
			r8 := e.r8Get(number)
			lower := r8 & 0xF
			upper := r8 >> 4
			result := (lower << 4) | upper
			e.r8Set(number, result)
			e.cpu.SetZeroFlag(result == 0)
			e.cpu.SetSubtractFlag(false)
			e.cpu.SetHalfCarryFlag(false)
			e.cpu.SetCarryFlag(false)
		case 56, 57, 58, 59, 60, 61, 62, 63: // SRL
			number := opcode & 0x7
			r8 := e.r8Get(number)
			bit0 := GetBit(r8, 0)
			result := r8 >> 1
			e.r8Set(number, result)
			e.cpu.SetZeroFlag(result == 0)
			e.cpu.SetSubtractFlag(false)
			e.cpu.SetHalfCarryFlag(false)
			e.cpu.SetCarryFlag(bit0)
		case 64, 65, 66, 67, 68, 69, 70, 71, // BIT bit, r8
			72, 73, 74, 75, 76, 77, 78, 79,
			80, 81, 82, 83, 84, 85, 86, 87,
			88, 89, 90, 91, 92, 93, 94, 95,
			96, 97, 98, 99, 100, 101, 102, 103,
			104, 105, 106, 107, 108, 109, 110, 111,
			112, 113, 114, 115, 116, 117, 118, 119,
			120, 121, 122, 123, 124, 125, 126, 127:
			bit := (opcode >> 3) & 0x7
			number := opcode & 0x7
			r8 := e.r8Get(number)
			e.cpu.SetSubtractFlag(false)
			e.cpu.SetHalfCarryFlag(true)
			e.cpu.SetZeroFlag(!GetBit(r8, int(bit)))
		case 128, 129, 130, 131, 132, 133, 134, 135, // RES bit, r8
			136, 137, 138, 139, 140, 141, 142, 143,
			144, 145, 146, 147, 148, 149, 150, 151,
			152, 153, 154, 155, 156, 157, 158, 159,
			160, 161, 162, 163, 164, 165, 166, 167,
			168, 169, 170, 171, 172, 173, 174, 175,
			176, 177, 178, 179, 180, 181, 182, 183,
			184, 185, 186, 187, 188, 189, 190, 191:
			bit := (opcode >> 3) & 0x7
			number := opcode & 0x7
			r8 := e.r8Get(number)
			e.r8Set(number, SetBit8(r8, bit, false))
		case 192, 193, 194, 195, 196, 197, 198, 199, // SET bit, r8
			200, 201, 202, 203, 204, 205, 206, 207,
			208, 209, 210, 211, 212, 213, 214, 215,
			216, 217, 218, 219, 220, 221, 222, 223,
			224, 225, 226, 227, 228, 229, 230, 231,
			232, 233, 234, 235, 236, 237, 238, 239,
			240, 241, 242, 243, 244, 245, 246, 247,
			248, 249, 250, 251, 252, 253, 254, 255:
			bit := (opcode >> 3) & 0x7
			number := opcode & 0x7
			r8 := e.r8Get(number)
			e.r8Set(number, SetBit8(r8, bit, true))
		default:
			log.Println("CB Opcode: ", opcode, " not found")
			return
		}
	default:
		log.Println("Opcode: ", opcode, " not found")
		return
	}

	if e.delayedActivateIMEatInstruction == e.numInstructions {
		e.IME = 1
	}
}
