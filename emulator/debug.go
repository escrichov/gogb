package emulator

import (
	"fmt"
	"strings"
)

func (e *Emulator) PrintRegisters() {
	registers := fmt.Sprintf("AF: %x (%s)\nBC: %x\nDE: %x\nHL: %x\nSP: %x\nPC: %x\n",
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
