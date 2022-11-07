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
	RomSize              uint8
	RamSize              uint8
	DestinationCode      uint8
	LicenseCodeOld       uint8
	MaskROMVersionNumber uint8
	ComplementCheck      uint8
	CheckSum             uint16
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
	e.romHeader.RomSize = e.rom0[0x148]
	e.romHeader.RamSize = e.rom0[0x149]
	e.romHeader.DestinationCode = e.rom0[0x14A]
	e.romHeader.LicenseCodeOld = e.rom0[0x14B]
	e.romHeader.MaskROMVersionNumber = e.rom0[0x14C]
	e.romHeader.ComplementCheck = e.rom0[0x14D]
	e.romHeader.CheckSum = binary.BigEndian.Uint16(e.rom0[0x14E:0x150])

	switch e.romHeader.CartridgeType {
	case 0:
		e.memoryBankController = 0
	case 1, 2, 3:
		e.memoryBankController = 1
		e.rom1Pointer = 1 << 14
	case 5, 6:
		e.memoryBankController = 2
	case 8, 9:
		e.memoryBankController = 0
	case 0xB, 0xC, 0xD:
		e.memoryBankController = 1
	case 0xF, 0x10, 0x11, 0x12, 0x13:
		e.memoryBankController = 3
		e.rom1Pointer = 32768
	case 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E:
		e.memoryBankController = 5
	}

	log.Println("Memory Bank Controller:", e.memoryBankController)

	return nil
}
