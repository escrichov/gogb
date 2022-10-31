package cpu

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
