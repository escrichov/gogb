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
		e.mem.GetDIV(), e.mem.io[263], e.mem.GetTIMA(), e.mem.GetTMA(), e.timer.internalTimer,
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

	var ramString string
	if mbcFeatures.RamSize > 0 {
		ramString = fmt.Sprintf("Ram Size: %d bytes (%s, %s) (%d Banks, 8KiB each)",
			mbcFeatures.RamSize,
			convertToHumanBytes(mbcFeatures.RamSize),
			convertToHumanBits(mbcFeatures.RamSize),
			mbcFeatures.RamBanks,
		)

		if mbcFeatures.RamFilename != "" {
			ramString += fmt.Sprintf(" - Ram Filename: %s", mbcFeatures.RamFilename)
		}
	} else {
		ramString = "No Ram"
	}

	var cgbString string
	if e.rom.features.SupportColor && e.rom.features.SupportMonochrome {
		cgbString = "Works in Monochrome and Gameboy Color"
	} else if e.rom.features.PGBMode {
		cgbString = "Works in special a non-CGB-mode called PGBMode"
	} else if e.rom.features.SupportMonochrome {
		cgbString = "Works in Monochrome only"
	} else if e.rom.features.SupportColor {
		cgbString = "Works in Gameboy Color only"
	}

	var sgbString string
	if e.rom.features.SupportSGB {
		sgbString = "Works on Super Gameboy"
	} else {
		sgbString = "Doesn't work on Super Gameboy"
	}

	cartridge := fmt.Sprintf(
		"Cartridge\n\t"+
			"Title: %s\n\t"+
			"Cartridge Type: %d %s\n\t"+
			"Rom Size: %d bytes (%s, %s) (%d Banks, 16KiB each)\n\t"+
			"%s\n\t"+
			"Battery: %t\n\t"+
			"License Code: %x, %s\n\t"+
			"Destination code: %x, %s\n\t"+
			"Mask ROM version number: %x\n\t"+
			"Manufacturer code: %s\n\t"+
			"%s (CGB Flag: %x)\n\t"+
			"%s (SGB flag: %x)\n\t"+
			"Header Checksum: %x (Ok: %t)\n\t"+
			"Global Checksum: %x (Ok: %t)\n\t"+
			"Logo ok: %t\n",
		e.rom.features.Title,
		e.rom.features.CartridgeType,
		e.rom.features.CartridgeTypeName,
		mbcFeatures.RomSize,
		convertToHumanBytes(mbcFeatures.RomSize),
		convertToHumanBits(mbcFeatures.RomSize),
		mbcFeatures.RomBanks,
		ramString,
		mbcFeatures.HasBattery,
		e.rom.features.LicenseCode,
		e.rom.features.LicenseCodeName,
		e.rom.features.DestinationCode,
		e.rom.features.DestinationCodeName,
		e.rom.features.MaskROMVersionNumber,
		e.rom.features.ManufacturerCode,
		cgbString,
		e.rom.features.ColorGB,
		sgbString,
		e.rom.features.GBSGBIndicator,
		e.rom.features.GlobalChecksum,
		e.rom.features.GlobalChecksumOk,
		e.rom.features.HeaderChecksum,
		e.rom.features.HeaderChecksumOk,
		e.rom.features.LogoOk,
	)
	fmt.Println(cartridge)
}
