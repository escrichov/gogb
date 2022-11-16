package emulator

type CPU struct {
	AF Register // Accumulator & Flags Register (ZNHC---) -> N & H flags are not used -> (Z--C---)
	BC Register
	DE Register
	HL Register

	SP Register
	PC uint16
}

type Register struct {
	// The value of the register.
	value uint16
}

// GetHi gets the higher byte of the register.
func (reg *Register) GetHi() byte {
	return byte(reg.value >> 8)
}

// GetLo gets the lower byte of the register.
func (reg *Register) GetLo() byte {
	return byte(reg.value & 0xFF)
}

// Get gets the 2 byte value of the register.
func (reg *Register) Get() uint16 {
	return reg.value
}

// SetHi sets the higher byte of the register.
func (reg *Register) SetHi(val byte) {
	reg.value = uint16(val)<<8 | (uint16(reg.value) & 0xFF)
}

// SetLo sets the lower byte of the register.
func (reg *Register) SetLo(val byte) {
	reg.value = uint16(val) | (uint16(reg.value) & 0xFF00)
}

// Set the value of the register.
func (reg *Register) Set(val uint16) {
	reg.value = val
}

// SetBit the value of the register.
func (reg *Register) SetBit(pos int, value bool) {
	SetBit16(&reg.value, pos, value)
}

// GetBit the value of the register.
func (reg *Register) GetBit(pos int) bool {
	return GetBit16(reg.value, pos)
}

// SetHighBit the value of the high part of the register.
func (reg *Register) SetHighBit(highPos int, value bool) {
	pos := highPos + 8
	reg.SetBit(pos, value)
}

// SetLowBit the value of the low part of the register.
func (reg *Register) SetLowBit(lowPos int, value bool) {
	reg.SetBit(lowPos, value)
}

// GetHighBit the value of the high part of the register.
func (reg *Register) GetHighBit(highPos int) bool {
	pos := highPos + 8
	return reg.GetBit(pos)
}

// GetLowBit the value of the low part of the register.
func (reg *Register) GetLowBit(lowPos int) bool {
	return reg.GetBit(lowPos)
}

// SetZeroFlag sets the value of the Zero flag.
func (cpu *CPU) SetZeroFlag(value bool) {
	cpu.AF.SetLowBit(7, value)
}

// SetSubtractFlag sets the value of Subtract flag.
func (cpu *CPU) SetSubtractFlag(value bool) {
	cpu.AF.SetLowBit(6, value)
}

// SetHalfCarryFlag sets the value of the Half Carry flag.
func (cpu *CPU) SetHalfCarryFlag(value bool) {
	cpu.AF.SetLowBit(5, value)
}

// SetCarryFlag sets the value of the Carry flag.
func (cpu *CPU) SetCarryFlag(value bool) {
	cpu.AF.SetLowBit(4, value)
}

// GetZeroFlag gets the value of the Zero flag.
func (cpu *CPU) GetZeroFlag() bool {
	return cpu.AF.GetLowBit(7)
}

// GetSubtractFlag gets the value of Subtract flag.
func (cpu *CPU) GetSubtractFlag() bool {
	return cpu.AF.GetLowBit(6)
}

// GetHalfCarryFlag gets the value of the Half Carry flag.
func (cpu *CPU) GetHalfCarryFlag() bool {
	return cpu.AF.GetLowBit(5)
}

// GetCarryFlag gets the value of the Carry flag.
func (cpu *CPU) GetCarryFlag() bool {
	return cpu.AF.GetLowBit(4)
}

func (cpu *CPU) GetA() uint8 {
	return cpu.AF.GetHi()
}

func (cpu *CPU) SetA(value uint8) {
	cpu.AF.SetHi(value)
}

func (cpu *CPU) GetF() uint8 {
	return cpu.AF.GetLo()
}

func (cpu *CPU) SetF(value uint8) {
	value &= 0xF0
	cpu.AF.SetLo(value)
}

func (cpu *CPU) GetB() uint8 {
	return cpu.BC.GetHi()
}

func (cpu *CPU) SetB(value uint8) {
	cpu.BC.SetHi(value)
}

func (cpu *CPU) GetC() uint8 {
	return cpu.BC.GetLo()
}

func (cpu *CPU) SetC(value uint8) {
	cpu.BC.SetLo(value)
}

func (cpu *CPU) GetD() uint8 {
	return cpu.DE.GetHi()
}

func (cpu *CPU) SetD(value uint8) {
	cpu.DE.SetHi(value)
}

func (cpu *CPU) GetE() uint8 {
	return cpu.DE.GetLo()
}

func (cpu *CPU) SetE(value uint8) {
	cpu.DE.SetLo(value)
}

func (cpu *CPU) GetH() uint8 {
	return cpu.HL.GetHi()
}

func (cpu *CPU) SetH(value uint8) {
	cpu.HL.SetHi(value)
}

func (cpu *CPU) GetL() uint8 {
	return cpu.HL.GetLo()
}

func (cpu *CPU) SetL(value uint8) {
	cpu.HL.SetLo(value)
}

func (cpu *CPU) GetAF() uint16 {
	return cpu.AF.Get()
}

func (cpu *CPU) SetAF(value uint16) {
	value &= 0xFFF0
	cpu.AF.Set(value)
}

func (cpu *CPU) GetBC() uint16 {
	return cpu.BC.Get()
}

func (cpu *CPU) SetBC(value uint16) {
	cpu.BC.Set(value)
}

func (cpu *CPU) GetDE() uint16 {
	return cpu.DE.Get()
}

func (cpu *CPU) SetDE(value uint16) {
	cpu.DE.Set(value)
}

func (cpu *CPU) GetHL() uint16 {
	return cpu.HL.Get()
}

func (cpu *CPU) SetHL(value uint16) {
	cpu.HL.Set(value)
}

func (cpu *CPU) GetSP() uint16 {
	return cpu.SP.Get()
}

func (cpu *CPU) SetSP(value uint16) {
	cpu.SP.Set(value)
}

func (cpu *CPU) checkCondition(conditionNumber uint8) bool {
	switch conditionNumber {
	case 0:
		return !cpu.GetZeroFlag()
	case 1:
		return cpu.GetZeroFlag()
	case 2:
		return !cpu.GetCarryFlag()
	case 3:
		return cpu.GetCarryFlag()
	default:
		return false
	}
}

func (cpu *CPU) r16Group1Get(number uint8) uint16 {
	switch number {
	case 0:
		return cpu.GetBC()
	case 1:
		return cpu.GetDE()
	case 2:
		return cpu.GetHL()
	case 3:
		return cpu.GetSP()
	default:
		return 0
	}
}

func (cpu *CPU) r16Group1Set(number uint8, val uint16) {
	switch number {
	case 0:
		cpu.SetBC(val)
	case 1:
		cpu.SetDE(val)
	case 2:
		cpu.SetHL(val)
	case 3:
		cpu.SetSP(val)
	default:
	}
}

func (cpu *CPU) r16Group2Get(number uint8) uint16 {
	switch number {
	case 0:
		return cpu.GetBC()
	case 1:
		return cpu.GetDE()
	case 2:
		value := cpu.GetHL()
		cpu.SetHL(value + 1)
		return value
	case 3:
		value := cpu.GetHL()
		cpu.SetHL(value - 1)
		return value
	default:
		return 0
	}
}

func (cpu *CPU) r16Group3Get(number uint8) uint16 {
	switch number {
	case 0:
		return cpu.GetBC()
	case 1:
		return cpu.GetDE()
	case 2:
		return cpu.GetHL()
	case 3:
		return cpu.GetAF()
	default:
		return 0
	}
}

func (cpu *CPU) r16Group3Set(number uint8, val uint16) {
	switch number {
	case 0:
		cpu.SetBC(val)
	case 1:
		cpu.SetDE(val)
	case 2:
		cpu.SetHL(val)
	case 3:
		cpu.SetAF(val)
	default:
	}
}

func (e *Emulator) popPC() uint8 {
	result := e.mem.read8(e.cpu.PC)
	e.cpu.PC++

	if e.isHaltBugActive {
		e.cpu.PC--
		e.isHaltBugActive = false
	}

	return result
}

func (e *Emulator) popPC16() uint16 {
	low := e.popPC()
	high := e.popPC()
	result := uint16(high)<<8 | uint16(low)
	return result
}

func (e *Emulator) r8Get(number uint8) uint8 {
	switch number {
	case 0:
		return e.cpu.GetB()
	case 1:
		return e.cpu.GetC()
	case 2:
		return e.cpu.GetD()
	case 3:
		return e.cpu.GetE()
	case 4:
		return e.cpu.GetH()
	case 5:
		return e.cpu.GetL()
	case 6:
		return e.mem.read8(e.cpu.GetHL())
	case 7:
		return e.cpu.GetA()
	default:
		return 0
	}
}

func (e *Emulator) r8Set(number uint8, val uint8) {
	switch number {
	case 0:
		e.cpu.SetB(val)
	case 1:
		e.cpu.SetC(val)
	case 2:
		e.cpu.SetD(val)
	case 3:
		e.cpu.SetE(val)
	case 4:
		e.cpu.SetH(val)
	case 5:
		e.cpu.SetL(val)
	case 6:
		e.mem.write8(e.cpu.GetHL(), val)
	case 7:
		e.cpu.SetA(val)
	default:
	}
}

func (e *Emulator) tick() {
	e.cycles += 4

	if e.timer.timaUpdateWithTMADelayedCycles == e.cycles {
		e.timer.reloadTIMAwithTMA()
	}
}

func (e *Emulator) push(val uint16) {
	sp := e.cpu.GetSP()
	sp--
	e.mem.write8(sp, uint8(val>>8))
	sp--
	e.mem.write8(sp, uint8(val))
	e.cpu.SetSP(sp)

	e.tick()
}

func (e *Emulator) pop() uint16 {
	sp := e.cpu.GetSP()
	result := e.mem.read16(sp)
	e.cpu.SetSP(sp + 2)

	return result
}
