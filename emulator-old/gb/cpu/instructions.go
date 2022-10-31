package cpu

var Instructions = [0x100]func(cpu *CPU){
	0x00: func(cpu *CPU) {
		// NOP
	},
	0x01: func(cpu *CPU) {
		// LD BC,u16
		u16 := popPC16(cpu)
		cpu.BC.Set(u16)
	},
	0x04: func(cpu *CPU) {
		// INC B
		result := cpu.BC.GetHi() + 1
		cpu.BC.SetHi(result)
		cpu.SetZeroFlag(result == 0)
		cpu.SetSubtractFlag(false)
		cpu.SetHalfCarryFlag(HalfCarryAdd(cpu.BC.GetHi(), 1))
	},
	0x05: func(cpu *CPU) {
		// DEC B
		b := cpu.BC.GetHi()
		result := b - 1
		cpu.BC.SetHi(result)
		cpu.SetZeroFlag(result == 0)
		cpu.SetSubtractFlag(true)
		cpu.SetHalfCarryFlag(b&0x0F == 0)
	},
	0x06: func(cpu *CPU) {
		// LD B,u8
		u8 := popPC(cpu)
		cpu.BC.SetHi(u8)
	},
	0x0C: func(cpu *CPU) {
		// INC C
		result := cpu.BC.GetLo() + 1
		cpu.BC.SetLo(result)
		cpu.SetZeroFlag(result == 0)
		cpu.SetSubtractFlag(false)
		cpu.SetHalfCarryFlag(HalfCarryAdd(cpu.BC.GetLo(), 1))
	},
	0x0D: func(cpu *CPU) {
		// DEC C
		c := cpu.BC.GetLo()
		result := c - 1
		cpu.BC.SetLo(result)
		cpu.SetZeroFlag(result == 0)
		cpu.SetSubtractFlag(true)
		cpu.SetHalfCarryFlag(c&0x0F == 0)
	},
	0x0E: func(cpu *CPU) {
		// LD C, u8
		u8 := popPC(cpu)
		cpu.BC.SetLo(u8)
	},
	0x11: func(cpu *CPU) {
		// LD DE,u16
		u16 := popPC16(cpu)
		cpu.DE.Set(u16)
	},
	0x13: func(cpu *CPU) {
		// INC DE
		cpu.DE.Set(cpu.DE.Get() + 1)
	},
	0x17: func(cpu *CPU) {
		// RLA
		a := cpu.AF.GetHi()
		oldCarryFlag := a >> 7
		carryFlag := byte(0)
		if cpu.GetCarryFlag() {
			carryFlag = byte(1)
		}
		rotatedA := byte(a<<1) | carryFlag
		cpu.AF.SetHi(rotatedA)
		cpu.SetZeroFlag(false)
		cpu.SetHalfCarryFlag(false)
		cpu.SetSubtractFlag(false)
		cpu.SetCarryFlag(oldCarryFlag == 1)
	},
	0x18: func(cpu *CPU) {
		// JR i8
		i8 := int8(popPC(cpu))
		addr := int32(cpu.PC) + int32(i8)
		cpu.instJump(uint16(addr))
	},
	0x1A: func(cpu *CPU) {
		// LD A, (DE)
		cpu.AF.SetHi(cpu.MMU.Read(cpu.DE.Get()))
	},
	0x1E: func(cpu *CPU) {
		// LD E,u8
		u8 := popPC(cpu)
		cpu.DE.SetLo(u8)
	},
	0x20: func(cpu *CPU) {
		// JR NZ,i8
		i8 := int8(popPC(cpu))
		if !cpu.GetZeroFlag() {
			addr := int32(cpu.PC) + int32(i8)
			cpu.instJump(uint16(addr))
			cpu.currentCpuCycles += 1
		}
	},
	0x21: func(cpu *CPU) {
		// LD HL, u16
		u16 := popPC16(cpu)
		cpu.HL.Set(u16)
	},
	0x22: func(cpu *CPU) {
		// LD (HL+),A
		addr := cpu.HL.Get()
		cpu.MMU.Write(addr, cpu.AF.GetHi())
		cpu.HL.Set(addr + 1)
	},
	0x23: func(cpu *CPU) {
		// INC HL
		cpu.HL.Set(cpu.HL.Get() + 1)
	},
	0x28: func(cpu *CPU) {
		// JR Z,i8
		i8 := int8(popPC(cpu))
		if cpu.GetZeroFlag() {
			addr := int32(cpu.PC) + int32(i8)
			cpu.instJump(uint16(addr))
			cpu.currentCpuCycles += 1
		}
	},
	0x2E: func(cpu *CPU) {
		// LD L,u8
		u8 := popPC(cpu)
		cpu.HL.SetLo(u8)
	},
	0x31: func(cpu *CPU) {
		// LD SP,u16
		u16 := popPC16(cpu)
		cpu.SP.Set(u16)
	},
	0x32: func(cpu *CPU) {
		// LD (HL-),A
		cpu.MMU.Write(cpu.HL.Get(), cpu.AF.GetHi())
		cpu.HL.Set(cpu.HL.Get() - 1)
	},
	0x3D: func(cpu *CPU) {
		// DEC A
		a := cpu.AF.GetHi()
		result := a - 1
		cpu.AF.SetHi(result)
		cpu.SetZeroFlag(result == 0)
		cpu.SetSubtractFlag(true)
		cpu.SetHalfCarryFlag(a&0x0F == 0)
	},
	0x3E: func(cpu *CPU) {
		// LD A,u8
		u8 := popPC(cpu)
		cpu.AF.SetHi(u8)
	},
	0x4F: func(cpu *CPU) {
		// LD C,A
		cpu.BC.SetLo(cpu.AF.GetHi())
	},
	0x57: func(cpu *CPU) {
		// LD D,A
		cpu.DE.SetHi(cpu.AF.GetHi())
	},
	0x6e: func(cpu *CPU) {
		// LD L, (HL)
		cpu.HL.SetLo(cpu.MMU.Read(cpu.HL.Get()))
	},
	0x67: func(cpu *CPU) {
		// LD H,A
		cpu.HL.SetHi(cpu.AF.GetHi())
	},
	0x77: func(cpu *CPU) {
		// LD (HL),A
		cpu.MMU.Write(cpu.HL.Get(), cpu.AF.GetHi())
	},
	0x7B: func(cpu *CPU) {
		// LD A,E
		cpu.AF.SetHi(cpu.DE.GetLo())
	},
	0x99: func(cpu *CPU) {
		// SBC A,C
		carry := 0
		if cpu.GetCarryFlag() {
			carry = 1
		}
		c := cpu.BC.GetLo()
		a := cpu.AF.GetHi()
		result16 := int16(a) - int16(c) - int16(carry)
		result := byte(result16)
		cpu.AF.SetHi(result)
		cpu.SetZeroFlag(result == 0)
		cpu.SetHalfCarryFlag(int16(a&0x0f)-int16(c&0xF)-int16(carry) < 0)
		cpu.SetSubtractFlag(true)
		cpu.SetCarryFlag(result16 < 0)
	},
	0xAF: func(cpu *CPU) {
		// XOR A,A
		a := cpu.AF.GetHi()
		result := a ^ a
		cpu.AF.SetHi(result)
		cpu.SetZeroFlag(result == 0)
		cpu.SetSubtractFlag(false)
		cpu.SetHalfCarryFlag(false)
		cpu.SetCarryFlag(false)
	},
	0xC1: func(cpu *CPU) {
		// POP BC
		result := cpu.insPop16()
		cpu.BC.Set(result)
	},
	0xC3: func(cpu *CPU) {
		// JP u16
		u16 := popPC16(cpu)
		cpu.instJump(u16)
	},
	0xC5: func(cpu *CPU) {
		// PUSH BC
		cpu.insPush16(cpu.BC.Get())
	},
	0xC9: func(cpu *CPU) {
		// RET
		result := cpu.insPop16()
		cpu.instJump(result)
	},
	0xCD: func(cpu *CPU) {
		// CALL u16
		u16 := popPC16(cpu)
		cpu.insPush16(cpu.PC)
		cpu.instJump(u16)
	},
	0xCC: func(cpu *CPU) {
		// CALL Z,u16
		u16 := popPC16(cpu)
		if cpu.GetZeroFlag() {
			cpu.insPush16(cpu.PC)
			cpu.instJump(u16)
		}
	},
	0xE0: func(cpu *CPU) {
		// LD (FF00+u8),A
		u8 := popPC(cpu)
		addr := 0xFF00 + uint16(u8)
		cpu.MMU.Write(addr, cpu.AF.GetHi())
	},
	0xE2: func(cpu *CPU) {
		// LD (FF00+C),A
		addr := 0xFF00 + uint16(cpu.BC.GetLo())
		cpu.MMU.Write(addr, cpu.AF.GetHi())
	},
	0xEA: func(cpu *CPU) {
		// LD (u16),A
		u16 := popPC16(cpu)
		cpu.MMU.Write(u16, cpu.AF.GetHi())
	},
	0xF0: func(cpu *CPU) {
		// LD A,(FF00+u8)
		u8 := popPC(cpu)
		addr := 0xFF00 + uint16(u8)
		cpu.AF.SetHi(cpu.MMU.Read(addr))
	},
	0xFE: func(cpu *CPU) {
		// CP A,u8
		u8 := popPC(cpu)
		a := cpu.AF.GetHi()
		result := a - u8
		cpu.SetZeroFlag(result == 0)
		cpu.SetHalfCarryFlag((u8 & 0x0f) > (a & 0x0f))
		cpu.SetSubtractFlag(true)
		cpu.SetCarryFlag(u8 > a)
	},
}

// TODO: Complete instructions
// TODO: Complete instructions CB

var InstructionsCB = [0x100]func(cpu *CPU){
	0x11: func(cpu *CPU) {
		// RL C
		c := cpu.BC.GetLo()
		carryBit := c >> 7
		rotatedC := (c << 1 & 0xFF) | carryBit
		cpu.BC.SetLo(rotatedC)
		cpu.SetZeroFlag(rotatedC == 0)
		cpu.SetHalfCarryFlag(false)
		cpu.SetSubtractFlag(false)
		cpu.SetCarryFlag(carryBit == 1)
	},
	0x7C: func(cpu *CPU) {
		// BIT 7,H
		bit := cpu.HL.GetHighBit(7)
		cpu.SetZeroFlag(!bit)
		cpu.SetHalfCarryFlag(true)
		cpu.SetSubtractFlag(false)
	},
}
