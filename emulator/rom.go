package emulator

import (
	"log"
	"os"
	"path/filepath"
)

type Rom struct {
	features   *RomFeatures
	controller MemoryBankController
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
