package mmu

import (
	"fmt"
)

type MMU struct {
	Memory [0xFFFF]byte
}

func (mmu *MMU) Init(bootRom []byte, gameRom []byte) error {
	if err := mmu.LoadRom(gameRom); err != nil {
		return err
	}
	if err := mmu.LoadBootRom(bootRom); err != nil {
		return err
	}

	return nil
}

func (mmu *MMU) Read(address uint16) byte {
	return mmu.Memory[address]
	switch {
	case address < 0x4000:
		// Cartridge ROM - Bank 0
		return 0x00
	case address < 0x8000:
		// Cartridge ROM - Banks 1-n
		return 0x00
	case address < 0xA000:
		// VRAM Banking
		return 0x00

	case address < 0xC000:
		// Cartridge RAM (External RAM)
		return 0x00

	case address < 0xD000:
		// Internal RAM - Bank 0
		return 0x00

	case address < 0xE000:
		// Internal RAM Bank 1-7
		return 0x00

	case address < 0xFE00:
		// Echo RAM
		// re-enable echo RAM?
		//mem.Data[address] = value
		//mem.Write(address-0x2000, value)
		return 0x00

	case address < 0xFEA0:
		// Object Attribute Memory (OAM)
		return 0x00

	case address < 0xFF00:
		// Unusable Memory
		return 0x00
	case address < 0xFF80:
		return mmu.readIORam(address)
	default:
		return mmu.readHighRam(address)
	}
}

// readIORam reads from 0xFF00-0xFF80 in the Memory address space. The range
// includes the hardware registers.
func (mmu *MMU) readIORam(address uint16) byte {
	switch {
	case address == 0xFF00:
		// Joypad address
		return 0x00
	case address >= 0xFF10 && address <= 0xFF26:
		// mem.gb.Sound.Read(address)
		return 0x00
	case address >= 0xFF30 && address <= 0xFF3F:
		// Writing to channel 3 waveform RAM.
		return 0x00
	case address == 0xFF0F:
		// mem.HighRAM[0x0F] | 0xE0
		return 0x00
	case address >= 0xFF72 && address <= 0xFF77:
		//log.Print("read from ", address)
		return 0x00
	case address == 0xFF68:
		// BG palette index
		return 0x00
	case address == 0xFF69:
		// BG Palette data
		return 0x00
	case address == 0xFF6A:
		// Sprite palette index
		return 0x00
	case address == 0xFF6B:
		// Sprite Palette data
		return 0x00
	case address == 0xFF4D:
		// Speed switch data
		return 0x00
	case address == 0xFF4F:
		// return mem.VRAMBank
		return 0x00
	case address == 0xFF70:
		// mem.WRAMBank
		return 0x00
	default:
		// mem.HighRAM[address-0xFF00]
		return 0x00
	}
}

// readHighRam reads from 0xFF80-0xFFFF in the Memory address space. The range
// includes HRAM.
func (mmu *MMU) readHighRam(address uint16) byte {
	switch {
	default:
		// mem.HighRAM[address-0xFF00]
		return 0x00
	}
}

func (mmu *MMU) Write(address uint16, value byte) {
	mmu.Memory[address] = value
	switch {
	case address < 0x4000:
		// Cartridge ROM - Bank 0
	case address < 0x8000:
		// Cartridge ROM - Banks 1-n
	case address < 0xA000:
		// VRAM Banking
	case address < 0xC000:
		// Cartridge RAM (External RAM)
	case address < 0xD000:
		// Internal RAM - Bank 0
	case address < 0xE000:
		// Internal RAM Bank 1-7
	case address < 0xFE00:
		// Echo RAM
		// re-enable echo RAM?
		//mem.Data[address] = value
		//mem.Write(address-0x2000, value)
	case address < 0xFEA0:
		// Object Attribute Memory (OAM)
	case address < 0xFF00:
		// Unusable Memory
	case address < 0xFF80:
		mmu.writeIORam(address, value)
	default:
		mmu.writeHighRam(address, value)
	}
}

// readIORam reads from 0xFF00-0xFF80 in the Memory address space. The range
// includes the hardware registers.
func (mmu *MMU) writeIORam(address uint16, value byte) {
	switch {
	case address == 0xFF00:
		// Joypad address
	case address >= 0xFF10 && address <= 0xFF26:
		// mem.gb.Sound.Read(address)
	case address >= 0xFF30 && address <= 0xFF3F:
		// Writing to channel 3 waveform RAM.
	case address == 0xFF0F:
		// mem.HighRAM[0x0F] | 0xE0
	case address >= 0xFF72 && address <= 0xFF77:
		//log.Print("read from ", address)
	case address == 0xFF68:
		// BG palette index
	case address == 0xFF69:
		// BG Palette data
	case address == 0xFF6A:
		// Sprite palette index
	case address == 0xFF6B:
		// Sprite Palette data
	case address == 0xFF4D:
		// Speed switch data
	case address == 0xFF4F:
		// return mem.VRAMBank
	case address == 0xFF70:
		// mem.WRAMBank
	}
}

// writeHighRam writes to 0xFF80-0xFFFF in the Memory address space. The range
// includes HRAM.
func (mmu *MMU) writeHighRam(address uint16, value byte) {
}

// TODO: DMA Transfer
// TODO: HDMA Transfer

func (mmu *MMU) LoadBootRom(data []byte) error {
	if len(data) != 0x100 {
		return fmt.Errorf("invalid length: %d", len(data))
	}

	copy(mmu.Memory[0:0xFF], data[0:0xFF])

	return nil
}

func (mmu *MMU) LoadRom(data []byte) error {
	if len(data) > 0x8000 {
		return fmt.Errorf("invalid length: %d", len(data))
	}

	copy(mmu.Memory[0x0:], data)

	return nil
}
