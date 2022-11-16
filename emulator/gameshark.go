package emulator

import (
	"encoding/binary"
	"fmt"
)

const (
	GamesharkCodePokemonBlueMissingNo          int = 0x0156D8CF
	GamesharkCodePokemonBlueNoRandomEncounters int = 0x01033CD1
	GamesharkCodePokemonBlueInfiniteMoney      int = 0x019947D3
	GamesharkCodePokemonBlueMasterballInMart   int = 0x01017CCF
	GamesharkCodePokemonBlueRareCandyInMart    int = 0x01287CCF
	GamesharkCodePokemonBlueHaveAll8Badge      int = 0x01FF56D3
	GamesharkCodePokemonBlueWalkThroughWalls   int = 0x010138CD
)

type GameSharkCode struct {
	address uint16
	value   uint8
}

func GamesharkParseCode(codeInt int) (*GameSharkCode, error) {
	// Game Shark codes consist of eight-digit hex numbers,
	// formatted as ABCDEFGH, the meaning of the separate digits is:
	// - AB External RAM bank number
	// - CD New Data
	// - GHEF Memory Address (internal or external RAM, A000-DFFF)

	codeType := (codeInt >> 24) & 0xFF
	if codeType != 1 {
		return nil, fmt.Errorf("unsupported code type %d", codeType)
	}

	var bs [2]byte
	var code GameSharkCode

	binary.LittleEndian.PutUint16(bs[:], uint16(codeInt))
	code.address = binary.BigEndian.Uint16(bs[:])
	code.value = uint8(codeInt >> 16)

	return &code, nil
}

func (e *Emulator) GamesharkAddCode(codeInt int) error {
	code, err := GamesharkParseCode(codeInt)
	if err != nil {
		return err
	}
	e.GameSharkCodes = append(e.GameSharkCodes, code)

	return nil
}

func (e *Emulator) loadGameSharkCodes() {
	for _, code := range e.GameSharkCodes {
		e.mem.write8(code.address, code.value)
	}
}
