package cpu

import "emulator-go/emulator/gb/utils"

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
	utils.SetBit16(&reg.value, pos, value)
}

// GetBit the value of the register.
func (reg *Register) GetBit(pos int) bool {
	return utils.GetBit16(reg.value, pos)
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
