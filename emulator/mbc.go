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
}

type BaseMBC struct {
	rom []byte
	ram []byte
	*MBCFeatures
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

func getMemoryBankControllerFeatures(romData []byte, romFilename string) (*MBCFeatures, error) {
	var mbcFeatures MBCFeatures

	cartridgeType := romData[0x147]
	romSizeByte := romData[0x148]
	ramSizeByte := romData[0x149]

	// ROM Size
	mbcFeatures.RomBankSize = 16384
	mbcFeatures.RomSize = 32768 * (1 << romSizeByte)
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

	// Check real Rom Size is equal to Rom header size
	if mbcFeatures.RomSize != len(romData) {
		return nil, fmt.Errorf(
			"real rom size (%d) != rom header size (%d)",
			mbcFeatures.RomSize,
			len(romData),
		)
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
	default:
		return nil, fmt.Errorf("unsupported memory bank controller: %d", baseMBC.MemoryBankControllerNumber)
	}

	return mbc, nil
}
