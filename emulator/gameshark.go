package emulator

import "encoding/binary"

func (e *Emulator) loadGamesharkCodes() {
	//e.ActivateCode(0x0156D8CF) // MissingNo
	//e.ActivateCode(0x01033CD1) // No Random Encounters
	//e.ActivateCode(0x019947D3) // Infinite Money
	//e.ActivateCode(0x01017CCF) // Masterball in Mart
	//e.ActivateCode(0x01287CCF) // Rare candy in Mart
	//e.ActivateCode(0x01FF56D3) // Have all 8 badge
	//e.ActivateCode(0x010138CD) // Walk through walls
}

func (e *Emulator) ActivateCode(code int) {
	codeType := (code >> 24) & 0xFF
	if codeType == 0x01 {
		bs := make([]byte, 2)
		binary.LittleEndian.PutUint16(bs, uint16(code))
		address := binary.BigEndian.Uint16(bs)
		value := uint8(code >> 16)

		e.write8(address, value)
	}
}
