package cpu

import (
	mmu "emulator-go/emulator/gb/memory"
)

// CPU contains the registers used for program execution and
// provides methods for setting flags.
type CPU struct {
	AF Register // Accumulator & Flags Register (ZNHC---) -> N & H flags are not used -> (Z--C---)
	BC Register
	DE Register
	HL Register

	SP Register
	PC uint16

	Divider int

	currentCpuCycles int

	MMU *mmu.MMU
}

// Init CPU and its registers to the initial values.
func (cpu *CPU) Init(cgb bool, MMU *mmu.MMU) {
	cpu.PC = 0x00
	cpu.AF.Set(0x0000)
	cpu.BC.Set(0x0000)
	cpu.DE.Set(0x0000)
	cpu.HL.Set(0x0000)
	cpu.SP.Set(0x0000)

	cpu.MMU = MMU
}
