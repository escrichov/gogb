package emulator

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type RomFeatures struct {
	Title                string
	TitleBytes           []byte
	ColorGB              uint8
	LicenseCodeNew       string
	GBSGBIndicator       uint8
	CartridgeType        uint8
	RomSizeByte          uint8
	RamSizeByte          uint8
	DestinationCode      uint8
	LicenseCodeOld       uint8
	MaskROMVersionNumber uint8
	ComplementCheck      uint8
	CheckSumBytes        []byte
	CheckSum             uint16
	CartridgeTypeName    string

	Filename string
}

type Rom struct {
	features   *RomFeatures
	controller MemoryBankController
}

func getCartridgeTypeName(cartridgeType uint8) string {
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

func newRomFromBytes(bs []byte, romFilename string) (*Rom, error) {
	var err error
	var rom = Rom{}

	rom.features, err = parseRomHeader(bs)
	if err != nil {
		return nil, err
	}

	rom.controller, err = newMemoryBankController(bs, romFilename)
	if err != nil {
		return nil, err
	}

	return &rom, err
}

func newRomFromFile(fileName string) (*Rom, error) {
	var err error
	var romData []byte

	bs, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Rom file not found:", err)
		return nil, err
	}

	fileExtension := filepath.Ext(fileName)
	if fileExtension == ".zip" {
		romData, err = UnzipBytes(bs)
		if err != nil {
			return nil, err
		}
	} else {
		romData = bs
	}

	return newRomFromBytes(romData, fileName)
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

func getMemoryBankControllerByCartridgeType(cartridgeType uint8) int {
	memoryBankControllerNumber := -1
	switch cartridgeType {
	case 0: // ROM ONLY
		memoryBankControllerNumber = 0
	case 1: // MBC1
		memoryBankControllerNumber = 1
	case 2: // MBC1+RAM
		memoryBankControllerNumber = 1
	case 3: // MBC1+RAM+BATTERY
		memoryBankControllerNumber = 1
	case 5:
		memoryBankControllerNumber = 2
	case 6:
		memoryBankControllerNumber = 2
	case 8:
		memoryBankControllerNumber = 0
	case 9:
		memoryBankControllerNumber = 0
	case 0xB, 0xC:
		memoryBankControllerNumber = 1
	case 0x0D:
		memoryBankControllerNumber = 1
	case 0x11, 0x12:
		memoryBankControllerNumber = 3
	case 0x0F, 0x10, 0x13:
		memoryBankControllerNumber = 3
	case 0x19, 0x1A, 0x1C, 0x1D:
		memoryBankControllerNumber = 5
	case 0x1B, 0x1E:
		memoryBankControllerNumber = 5
	}

	return memoryBankControllerNumber
}

func hasBattery(cartridgeType uint8) bool {
	switch cartridgeType {
	case 0x03: // MBC1+RAM+BATTERY
		return true
	case 0x06: // MBC2+BATTERY
		return true
	case 0x09: // ROM+RAM+BATTERY
		return true
	case 0x0D: // MMM01+RAM+BATTERY
		return true
	case 0x0F: // MBC3+TIMER+BATTERY
		return true
	case 0x10: // MBC3+TIMER+RAM+BATTERY
		return true
	case 0x13: // MBC3+RAM+BATTERY
		return true
	case 0x1B: // MBC5+RAM+BATTERY
		return true
	case 0x1E: // MBC5+RUMBLE+RAM+BATTERY
		return true
	case 0x22: // MBC7+SENSOR+RUMBLE+RAM+BATTERY
		return true
	case 0xFF: // HuC1+RAM+BATTERY
		return true
	}

	return false
}

func isRamAllowed(cartridgeType uint8) bool {
	switch cartridgeType {
	case 0x02: // MBC1+RAM
		return true
	case 0x03: // MBC1+RAM+BATTERY
		return true
	case 0x05: // MBC2 (It always has ram)
		return true
	case 0x06: // MBC2+BATTERY (It always has ram)
		return true
	case 0x08: // ROM+RAM
		return true
	case 0x09: // ROM+RAM+BATTERY
		return true
	case 0x0C: // MMM01+RAM
		return true
	case 0x0D: // MMM01+RAM+BATTERY
		return true
	case 0x10: // MBC3+TIMER+RAM+BATTERY
		return true
	case 0x12: // MBC3+RAM
		return true
	case 0x13: // MBC3+RAM+BATTERY
		return true
	case 0x1A: // MBC5+RAM
		return true
	case 0x1B: // MBC5+RAM+BATTERY
		return true
	case 0x1D: // MBC5+RUMBLE+RAM
		return true
	case 0x1E: // MBC5+RUMBLE+RAM+BATTERY
		return true
	case 0x22: // MBC7+SENSOR+RUMBLE+RAM+BATTERY
		return true
	case 0xFF: // HuC1+RAM+BATTERY
		return true
	}

	return false
}

func parseRomHeader(romData []byte) (*RomFeatures, error) {
	var romFeatures RomFeatures

	if len(romData) < 0x150 {
		return nil, fmt.Errorf("incorrect rom size %d", len(romData))
	}

	romFeatures.Title = parseRomTitle(romData[0x134:0x144])
	romFeatures.TitleBytes = romData[0x134:0x144]
	romFeatures.ColorGB = romData[0x143]
	romFeatures.LicenseCodeNew = string(romData[0x144:0x146])
	romFeatures.GBSGBIndicator = romData[0x146]
	romFeatures.CartridgeType = romData[0x147]
	romFeatures.RomSizeByte = romData[0x148]
	romFeatures.RamSizeByte = romData[0x149]
	romFeatures.DestinationCode = romData[0x14A]
	romFeatures.LicenseCodeOld = romData[0x14B]
	romFeatures.MaskROMVersionNumber = romData[0x14C]
	romFeatures.ComplementCheck = romData[0x14D]
	romFeatures.CheckSumBytes = romData[0x14E:0x150]
	romFeatures.CheckSum = binary.BigEndian.Uint16(romFeatures.CheckSumBytes)
	romFeatures.CartridgeTypeName = getCartridgeTypeName(romFeatures.CartridgeType)

	return &romFeatures, nil
}
