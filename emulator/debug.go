package emulator

import (
	"fmt"
	"strings"
)

func (e *Emulator) PrintRegisters() {
	registers := fmt.Sprintf("Registers\n\tAF: %x (%s)\n\tBC: %x\n\tDE: %x\n\tHL: %x\n\tSP: %x\n\tPC: %x\n",
		e.cpu.AF.value,
		e.GetDebugFlags(),
		e.cpu.BC.value,
		e.cpu.DE.value,
		e.cpu.HL.value,
		e.cpu.SP.value,
		e.cpu.PC)
	fmt.Println(registers)
}

func (e *Emulator) GetDebugFlags() string {
	var flags strings.Builder
	if e.cpu.GetZeroFlag() {
		flags.WriteRune('Z')
	} else {
		flags.WriteRune('-')
	}

	if e.cpu.GetSubtractFlag() {
		flags.WriteRune('N')
	} else {
		flags.WriteRune('-')
	}

	if e.cpu.GetHalfCarryFlag() {
		flags.WriteRune('H')
	} else {
		flags.WriteRune('-')
	}

	if e.cpu.GetCarryFlag() {
		flags.WriteRune('C')
	} else {
		flags.WriteRune('-')
	}

	return flags.String()
}

func (e *Emulator) PrintTimers() {
	timers := fmt.Sprintf("Timers\n\tDIV: %x\n\tTAC: %x\n\tTIMA: %x\n\tTMA: %x\n\tInternal Timer: %x\n",
		e.GetDIV(), e.io[263], e.GetTIMA(), e.GetTMA(), e.internalTimer,
	)
	fmt.Println(timers)
}

func (e *Emulator) PrintCartridge() {
	cartridge := fmt.Sprintf(
		"Cartridge\n\t"+
			"Title: %s\n\t"+
			"Cartridge Type: %s (%d)\n\t"+
			"MBC: %d\n\t"+
			"Rom Size: %d\n\t"+
			"Rom Banks: %d (16KiB each)\n\t"+
			"Ram Size: %d\n\t"+
			"Ram Banks: %d (8KiB each)\n\t"+
			"CGB flag: %x\n",
		e.romHeader.Title,
		e.romHeader.CartridgeTypeName,
		e.romHeader.CartridgeType,
		e.memoryBankController,
		e.romHeader.RomSize,
		e.romHeader.RomBanks,
		e.romHeader.RamSize,
		e.romHeader.RamBanks,
		e.romHeader.ColorGB,
	)
	fmt.Println(cartridge)
}
