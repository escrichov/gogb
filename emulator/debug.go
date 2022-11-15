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

func convertToHumanBytes(numBytes int) string {
	if numBytes < 1024 {
		return fmt.Sprintf("%d Bytes", numBytes/1024)
	} else if numBytes >= 1024 && numBytes < 1048576 {
		return fmt.Sprintf("%d KiB", numBytes/1024)
	} else {
		return fmt.Sprintf("%d MiB", numBytes/1048576)
	}
}

func convertToHumanBits(numBytes int) string {
	numBits := numBytes * 8
	if numBits < 1024 {
		return fmt.Sprintf("%d bits", numBits/1024)
	} else if numBits >= 1024 && numBits < 1048576 {
		return fmt.Sprintf("%d Kib", numBits/1024)
	} else {
		return fmt.Sprintf("%d Mib", numBits/1048576)
	}
}

func (e *Emulator) PrintCartridge() {
	mbcFeatures := e.rom.controller.GetFeatures()

	cartridge := fmt.Sprintf(
		"Cartridge\n\t"+
			"Title: %s\n\t"+
			"Cartridge Type: %s (%d)\n\t"+
			"MBC: %d\n\t"+
			"Rom Size: %d bytes (%s, %s)\n\t"+
			"Rom Banks: %d (16KiB each)\n\t"+
			"Ram Size: %d bytes (%s, %s)\n\t"+
			"Ram Banks: %d (8KiB each)\n\t"+
			"Battery: %t\n\t"+
			"Ram Filename: %s\n\t"+
			"CGB flag: %x\n",
		e.rom.features.Title,
		e.rom.features.CartridgeTypeName,
		e.rom.features.CartridgeType,
		mbcFeatures.MemoryBankControllerNumber,
		mbcFeatures.RomSize,
		convertToHumanBytes(mbcFeatures.RomSize),
		convertToHumanBits(mbcFeatures.RomSize),
		mbcFeatures.RomBanks,
		mbcFeatures.RamSize,
		convertToHumanBytes(mbcFeatures.RamSize),
		convertToHumanBits(mbcFeatures.RamSize),
		mbcFeatures.RamBanks,
		mbcFeatures.HasBattery,
		mbcFeatures.RamFilename,
		e.rom.features.ColorGB,
	)
	fmt.Println(cartridge)
}
