package cpu

var instIndex uint64 = 0

// ExecuteNextOpcode gets the value at the current PC address, increments the PC,
// updates the CPU ticks and executes the opcode.
func ExecuteNextOpcode(cpu *CPU) int {

	opcode := popPC(cpu)
	if opcode == 0xCB {
		opcodeCB := popPC(cpu)
		//fmt.Printf("Opcode CB %x - PC: %x\n - Index: %d", opcodeCB, pc, instIndex)
		cpu.currentCpuCycles = CBOpcodeCycles[opcodeCB]
		InstructionsCB[opcodeCB](cpu)
	} else {
		//fmt.Printf("Opcode: %x - PC: %x - Index: %d\n", opcode, pc, instIndex)
		cpu.currentCpuCycles = OpcodeCycles[opcode]
		Instructions[opcode](cpu)
	}
	instIndex += 1

	return cpu.currentCpuCycles * 4
}

// Read the value at the PC and increment the PC.
func popPC(cpu *CPU) byte {
	opcode := cpu.MMU.Read(cpu.PC)
	cpu.PC++
	return opcode
}

// Read the next 16bit value at the PC.
func popPC16(cpu *CPU) uint16 {
	b1 := uint16(popPC(cpu))
	b2 := uint16(popPC(cpu))
	return b2<<8 | b1
}
