package emulator

import (
	"fmt"
	"path/filepath"
)

type MBCFeatures struct {
	MemoryBankControllerNumber int

	RomSize     int
	RomBanks    uint16
	RomBankSize int

	RamSize     int
	RamBanks    uint16
	RamBankSize int
	RamFilename string

	HasBattery bool
}

type MemoryBankController interface {
	Read(address uint16) uint8
	Write(address uint16, val uint8)
	GetFeatures() *MBCFeatures
	GetRam() []byte
}

type BaseMBC struct {
	rom []byte
	ram []byte
	*MBCFeatures
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
	case 0x20:
		memoryBankControllerNumber = 6
	case 0x22:
		memoryBankControllerNumber = 7
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

func newBaseMBC(romData []byte, romFilename string) (*BaseMBC, error) {
	var err error
	var mbc = BaseMBC{rom: romData}

	mbc.MBCFeatures, err = getMemoryBankControllerFeatures(romData, romFilename)
	if err != nil {
		return nil, err
	}

	if mbc.RamSize > 0 {
		if mbc.HasBattery {
			mbc.ram, err = createMMAP(mbc.RamFilename, mbc.RamSize)
			if err != nil {
				return nil, err
			}
		} else {
			mbc.ram = make([]byte, mbc.RamSize)
		}
	}

	return &mbc, nil
}

func (mbc *BaseMBC) GetFeatures() *MBCFeatures {
	return mbc.MBCFeatures
}

func (mbc *BaseMBC) GetRam() []byte {
	return mbc.ram
}

func getMemoryBankControllerFeatures(romData []byte, romFilename string) (*MBCFeatures, error) {
	var mbcFeatures MBCFeatures

	cartridgeType := romData[0x147]
	romSizeByte := romData[0x148]
	ramSizeByte := romData[0x149]

	// ROM Size
	mbcFeatures.RomBankSize = 16384
	mbcFeatures.RomSize = 32768 * (1 << romSizeByte)
	// Check real Rom Size is equal to Rom header size
	if mbcFeatures.RomSize != len(romData) {
		return nil, fmt.Errorf(
			"real rom size (%d bytes) != rom header size (%d bytes)",
			len(romData),
			mbcFeatures.RomSize,
		)
	}
	mbcFeatures.RomBanks = uint16(mbcFeatures.RomSize / mbcFeatures.RomBankSize)

	// RAM Size
	switch ramSizeByte {
	case 0, 1:
		mbcFeatures.RamSize = 0
	case 2:
		mbcFeatures.RamSize = 8192
	case 3:
		mbcFeatures.RamSize = 32768
	case 4:
		mbcFeatures.RamSize = 131072
	case 5:
		mbcFeatures.RamSize = 65536
	}

	// RAM Banks
	mbcFeatures.RamBankSize = 8192
	mbcFeatures.RamBanks = uint16(mbcFeatures.RamSize / mbcFeatures.RamBankSize)

	// Check Ram Size is allowed in the Cartridge Header
	if !isRamAllowed(cartridgeType) && mbcFeatures.RamSize > 0 {
		return nil, fmt.Errorf("ram is not allowed for this cartridge type")
	}

	// Battery
	mbcFeatures.HasBattery = hasBattery(cartridgeType)

	// Memory Bank Controller
	mbcFeatures.MemoryBankControllerNumber = getMemoryBankControllerByCartridgeType(cartridgeType)

	// Ram filename
	if mbcFeatures.RamSize > 0 {
		extension := filepath.Ext(romFilename)
		mbcFeatures.RamFilename = romFilename[:len(romFilename)-len(extension)] + ".sav"
	}

	return &mbcFeatures, nil
}

func newMemoryBankController(romData []byte, romFilename string) (MemoryBankController, error) {
	var mbc MemoryBankController

	baseMBC, err := newBaseMBC(romData, romFilename)
	if err != nil {
		return nil, err
	}

	switch baseMBC.MemoryBankControllerNumber {
	case 0:
		mbc = NewMBC0(baseMBC)
	case 1:
		mbc = NewMBC1(baseMBC)
	case 2:
		mbc = NewMBC2(baseMBC)
	case 3:
		mbc = NewMBC3(baseMBC)
	case 5:
		mbc = NewMBC5(baseMBC)
	case 6:
		mbc = NewMBC5(baseMBC)
	default:
		return nil, fmt.Errorf("unsupported memory bank controller: %d", baseMBC.MemoryBankControllerNumber)
	}

	return mbc, nil
}
