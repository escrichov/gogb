package cpu

// Perform a JUMP operation by setting the PC to the value.
func (cpu *CPU) instJump(next uint16) {
	cpu.PC = next
}

func (cpu *CPU) insPush16(addr uint16) {
	sp := cpu.SP.Get()
	cpu.MMU.Write(sp-1, byte(addr>>8))
	cpu.MMU.Write(sp-2, byte(addr&0xFF))
	cpu.SP.Set(sp - 2)
}

func (cpu *CPU) insPop16() uint16 {
	sp := cpu.SP.Get()
	var result = uint16(cpu.MMU.Read(sp))
	result |= uint16(cpu.MMU.Read(sp+1)) << 8
	cpu.SP.Set(sp + 2)

	return result
}
