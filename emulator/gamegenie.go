package emulator

import (
	"fmt"
	"strconv"
)

const (
	GameGeniePokemonBlueDebugMode                         string = "CED-56A-D50"
	GameGeniePokemonBlueItemsCost0                        string = "1A8-258-A22"
	GameGeniePokemonBlueFirstStarterPokeballMewtwoo       string = "831-0EA-2A8"
	GameGeniePokemonBlueFirstStarterPokeballMew           string = "151-0EA-2A8"
	GameGeniePokemonBlueFirstStarterPokeballPikachu       string = "541-0EA-2A8"
	GameGeniePokemonBlueFirstStarterPokeballEevee         string = "661-0EA-2A8"
	GameGeniePokemonBlueAlwaysEncounterWildMissingNo      string = "03F-389-F7A"
	GameGeniePokemonBlueAlwaysEncounterPureWildMissingNo  string = "00B-059-F7A"
	GameGeniePokemonBlueAlwaysCatchPokemonWithAnyPokeball string = "C37-04A-C41"
	GameGeniePokemonBlueFasterDialogue                    string = "009-0EC-7F5"
	GameGeniePokemonBlueInfinitePP                        string = "C90-2DB-3BE"
	GameGeniePokemonBlueNoRandomBattles                   string = "C9F-3C9-E69"
)

type GameGenieCode struct {
	address  uint16
	newValue uint8
	oldValue uint8
}

func hex2uint(hexStr string) (uint64, error) {
	// base 16 for hexadecimal
	return strconv.ParseUint(hexStr, 16, 64)
}

// GameGenieParseCode
// Game Genie codes consist of nine-digit hex numbers,
// formatted as ABC-DEF-GHI, the meaning of the separate digits is:
// - AB, new data
// - FCDE, memory address, XORed by 0F000h
// - GI, old data, XORed by 0BAh and rotated left by two
// - H, Unknown, maybe checksum and/or else
func GameGenieParseCode(codeString string) (*GameGenieCode, error) {
	if len(codeString) != 11 {
		return nil, fmt.Errorf("invalid GameGenie code format: %s", codeString)
	}

	if codeString[3] != '-' || codeString[7] != '-' {
		return nil, fmt.Errorf("invalid GameGenie code format: %s", codeString)
	}

	var code GameGenieCode

	// New Value
	newValue, err := hex2uint(codeString[0:2])
	if err != nil {
		return nil, err
	}
	code.newValue = uint8(newValue)

	// Address
	address := string(codeString[6]) + string(codeString[2]) + codeString[4:6]
	addressUint, err := hex2uint(address)
	if err != nil {
		return nil, err
	}
	code.address = uint16(addressUint) ^ 0x0F000

	// Old value
	oldValueString := string(codeString[8]) + string(codeString[10])
	oldValueUint, err := hex2uint(oldValueString)
	if err != nil {
		return nil, err
	}
	oldValue := uint8(oldValueUint<<6) | uint8(oldValueUint>>2) // Rotate 2 bits to the right
	code.oldValue = oldValue ^ 0xBA

	// The address should be located in ROM area 0000h-7FFFh
	if code.address > 0x7FFF {
		return nil, fmt.Errorf("invalid address: %x. It should be location in 0x0000-0x7FFF area", code.address)
	}

	return &code, nil
}

func (e *Emulator) GameGenieAddCode(codeString string) error {
	code, err := GameGenieParseCode(codeString)
	if err != nil {
		return err
	}
	e.GameGenieCodes = append(e.GameGenieCodes, code)

	return nil
}

func (e *Emulator) getGameGenieValue(address uint16, value uint8) uint8 {
	for _, code := range e.GameGenieCodes {
		if code.address == address {
			if code.oldValue == value {
				return code.newValue
			}
		}
	}

	return value
}
