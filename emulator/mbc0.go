package emulator

type MBC0 struct {
	*BaseMBC
}

func NewMBC0(mbc *BaseMBC) *MBC0 {
	return &MBC0{BaseMBC: mbc}
}

func (mbc *MBC0) Read(addr uint16) uint8 {
	switch addr >> 13 {
	case 0, 1, 2, 3: // 0x0000 - 0x7FFF
		return mbc.rom[addr]
	case 5: // 0xA000 - 0xBFFF
		return 0xFF
	}

	return 0xFF
}

func (mbc *MBC0) Write(addr uint16, val uint8) {
	// Not allowed
}
