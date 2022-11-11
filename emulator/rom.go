package emulator

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type RomHeader struct {
	Title                string
	TitleBytes           []byte
	ColorGB              uint8
	LicenseCodeNew       string
	GBSGBIndicator       uint8
	CartridgeType        uint8
	RomSizeByte          uint8
	RomBanks             uint16
	RomSize              int
	RamSizeByte          uint8
	RamSize              int
	RamBanks             uint16
	DestinationCode      uint8
	LicenseCodeOld       uint8
	MaskROMVersionNumber uint8
	ComplementCheck      uint8
	CheckSum             uint16
	CartridgeTypeName    string
	HasBattery           bool
}

func romName(cartridgeType uint8) string {
	switch cartridgeType {
	case 0x00:
		return "ROM ONLY"
	case 0x01:
		return "MBC1"
	case 0x02:
		return "MBC1+RAM"
	case 0x03:
		return "MBC1+RAM+BATTERY"
	case 0x05:
		return "MBC2"
	case 0x06:
		return "MBC2+BATTERY"
	case 0x08:
		return "ROM+RAM 1"
	case 0x09:
		return "ROM+RAM+BATTERY 1"
	case 0x0B:
		return "MMM01"
	case 0x0C:
		return "MMM01+RAM"
	case 0x0D:
		return "MMM01+RAM+BATTERY"
	case 0x0F:
		return "MBC3+TIMER+BATTERY"
	case 0x10:
		return "MBC3+TIMER+RAM+BATTERY 2"
	case 0x11:
		return "MBC3"
	case 0x12:
		return "MBC3+RAM 2"
	case 0x13:
		return "MBC3+RAM+BATTERY 2"
	case 0x19:
		return "MBC5"
	case 0x1A:
		return "MBC5+RAM"
	case 0x1B:
		return "MBC5+RAM+BATTERY"
	case 0x1C:
		return "MBC5+RUMBLE"
	case 0x1D:
		return "MBC5+RUMBLE+RAM"
	case 0x1E:
		return "MBC5+RUMBLE+RAM+BATTERY"
	case 0x20:
		return "MBC6"
	case 0x22:
		return "MBC7+SENSOR+RUMBLE+RAM+BATTERY"
	case 0xFC:
		return "POCKET CAMERA"
	case 0xFD:
		return "BANDAI TAMA5"
	case 0xFE:
		return "HuC3"
	case 0xFF:
		return "HuC1+RAM+BATTERY"
	default:
		return "Unknown"
	}
}

func (e *Emulator) loadRom(fileName string) error {
	var err error

	bs, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Rom file not found:", err)
		return err
	}

	fileExtension := filepath.Ext(fileName)
	if fileExtension == ".zip" {
		e.rom0, err = UnzipBytes(bs)
		if err != nil {
			return err
		}
	} else {
		e.rom0 = bs
	}

	err = e.parseRomHeader()
	if err != nil {
		return err
	}

	return nil
}

func parseRomTitle(bs []byte) string {
	var end = 0
	for n, b := range bs {
		end = n
		if b == 0 {
			break
		}
	}
	return string(bs[0:end])
}

func (e *Emulator) parseRomHeader() error {
	if len(e.rom0) < 0x150 {
		return fmt.Errorf("incorrect rom size %d", len(e.rom0))
	}

	e.romHeader.Title = parseRomTitle(e.rom0[0x134:0x144])
	e.romHeader.TitleBytes = e.rom0[0x134:0x144]
	e.romHeader.ColorGB = e.rom0[0x143]
	e.romHeader.LicenseCodeNew = string(e.rom0[0x144:0x146])
	e.romHeader.GBSGBIndicator = e.rom0[0x146]
	e.romHeader.CartridgeType = e.rom0[0x147]
	e.romHeader.RomSizeByte = e.rom0[0x148]
	e.romHeader.RamSizeByte = e.rom0[0x149]
	e.romHeader.DestinationCode = e.rom0[0x14A]
	e.romHeader.LicenseCodeOld = e.rom0[0x14B]
	e.romHeader.MaskROMVersionNumber = e.rom0[0x14C]
	e.romHeader.ComplementCheck = e.rom0[0x14D]
	e.romHeader.CheckSum = binary.BigEndian.Uint16(e.rom0[0x14E:0x150])

	e.romHeader.CartridgeTypeName = romName(e.romHeader.CartridgeType)

	// ROM
	e.romHeader.RomSize = 32768 * (1 << e.romHeader.RomSizeByte)
	e.romHeader.RomBanks = uint16(e.romHeader.RomSize / 16384)

	// Check real Rom Size is equal to Rom header size
	if e.romHeader.RomSize != len(e.rom0) {
		return fmt.Errorf("real rom size (%d) != rom header size (%d)", e.romHeader.RomSize, len(e.rom0))
	}

	// RAM
	switch e.romHeader.RamSizeByte {
	case 0, 1:
		e.romHeader.RamSize = 0
		e.romHeader.RamBanks = 0
	case 2:
		e.romHeader.RamSize = 8192
		e.romHeader.RamBanks = 1
	case 3:
		e.romHeader.RamSize = 32768
		e.romHeader.RamBanks = 4
	case 4:
		e.romHeader.RamSize = 131072
		e.romHeader.RamBanks = 16
	case 5:
		e.romHeader.RamSize = 65536
		e.romHeader.RamBanks = 8
	}

	switch e.romHeader.CartridgeType {
	case 0: // ROM ONLY
		e.memoryBankController = 0
	case 1: // MBC1
		e.memoryBankController = 1
		e.romHeader.RamSize = 0
		e.romHeader.RamBanks = 0
	case 2: // MBC1+RAM
		e.memoryBankController = 1
	case 3: // MBC1+RAM+BATTERY
		e.memoryBankController = 1
		e.romHeader.HasBattery = true
	case 5:
		e.memoryBankController = 2
	case 6:
		e.romHeader.HasBattery = true
		e.memoryBankController = 2
	case 8:
		e.memoryBankController = 0
	case 9:
		e.memoryBankController = 0
		e.romHeader.HasBattery = true
	case 0xB, 0xC:
		e.memoryBankController = 1
	case 0x0D:
		e.memoryBankController = 1
		e.romHeader.HasBattery = true
	case 0x11, 0x12:
		e.memoryBankController = 3
	case 0x0F, 0x10, 0x13:
		e.memoryBankController = 3
		e.romHeader.HasBattery = true
	case 0x19, 0x1A, 0x1C, 0x1D:
		e.memoryBankController = 5
	case 0x1B, 0x1E:
		e.memoryBankController = 5
		e.romHeader.HasBattery = true
	case 0x22:
		e.romHeader.HasBattery = true
	case 0xFF:
		e.romHeader.HasBattery = true
	}

	// Default values for memory bank controllers
	if e.memoryBankController == 1 {
		e.mbc1Bank1 = 1
		if e.romHeader.RomSize >= 1048576 {
			e.mbc1AllowedRomBank2 = true
		}
		if e.romHeader.RamSize >= 32768 {
			e.mbc1AllowedRamBank2 = true
		}
	} else if e.memoryBankController == 2 {
		e.mbc2RomBank = 1
	} else if e.memoryBankController == 3 {
		e.mbc3RomBank = 1
		e.mbc3LatchRegister = 0xFFFF
	} else if e.memoryBankController == 5 {
		e.mbc5RomBank = 1
	}

	e.PrintCartridge()

	return nil
}
