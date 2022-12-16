package emulator

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

func WriteMBCBlock(data *bytes.Buffer, controller MemoryBankController) {
	data.WriteString("MBC ")

	switch mbc := controller.(type) {
	case *MBC0:
		binary.Write(data, binary.LittleEndian, uint32(0))
	case *MBC1:
		binary.Write(data, binary.LittleEndian, uint32(12))

		// 0000–1FFF — RAM Enable (Write Only)
		binary.Write(data, binary.LittleEndian, uint16(0x0000))
		if mbc.enableRamBank {
			binary.Write(data, binary.LittleEndian, uint8(0x0A))
		} else {
			binary.Write(data, binary.LittleEndian, uint8(0))
		}

		// 2000–3FFF — ROM Bank Number (Write Only)
		binary.Write(data, binary.LittleEndian, uint16(0x2000))
		binary.Write(data, binary.LittleEndian, mbc.bank1)

		// 4000–5FFF — RAM Bank Number — or — Upper Bits of ROM Bank Number (Write Only)
		binary.Write(data, binary.LittleEndian, uint16(0x4000))
		binary.Write(data, binary.LittleEndian, mbc.bank2)

		// 6000–7FFF — Banking Mode Select (Write Only)
		binary.Write(data, binary.LittleEndian, uint16(0x6000))
		if mbc.memoryModel == 0 {
			binary.Write(data, binary.LittleEndian, uint8(0x0))
		} else {
			binary.Write(data, binary.LittleEndian, uint8(0x1))
		}
	case *MBC2:
		binary.Write(data, binary.LittleEndian, uint32(6))

		// 0000–3FFF — RAM Enable
		binary.Write(data, binary.LittleEndian, uint16(0x0000))
		if mbc.enableRamBank {
			binary.Write(data, binary.LittleEndian, uint8(0x0A))
		} else {
			binary.Write(data, binary.LittleEndian, uint8(0x0))
		}

		// 0000–3FFF — ROM Bank Number
		binary.Write(data, binary.LittleEndian, uint16(0x0000))
		binary.Write(data, binary.LittleEndian, 0x10|mbc.romBank)
	case *MBC3:
		binary.Write(data, binary.LittleEndian, uint32(15))

		// A000-BFFF - RTC Register 08-0C
		binary.Write(data, binary.LittleEndian, uint16(0xA000))
		binary.Write(data, binary.LittleEndian, uint8(0))

		// 0000-1FFF - RAM and Timer Enable
		binary.Write(data, binary.LittleEndian, uint16(0x0000))
		if mbc.enableRamBank {
			binary.Write(data, binary.LittleEndian, uint8(0x0A))
		} else {
			binary.Write(data, binary.LittleEndian, uint8(0x0))
		}

		// 2000-3FFF - ROM Bank Number
		binary.Write(data, binary.LittleEndian, uint16(0x2000))
		binary.Write(data, binary.LittleEndian, mbc.romBank)

		// 4000-5FFF - RAM Bank Number - or - RTC Register Select
		binary.Write(data, binary.LittleEndian, uint16(0x4000))
		binary.Write(data, binary.LittleEndian, mbc.ramBank)

		// 6000-7FFF - Latch Clock Data
		binary.Write(data, binary.LittleEndian, uint16(0x6000))
		binary.Write(data, binary.LittleEndian, uint8(0))
	case *MBC5:
		binary.Write(data, binary.LittleEndian, uint32(12))

		// 0000-1FFF - RAM Enable
		binary.Write(data, binary.LittleEndian, uint16(0x0000))
		if mbc.enableRamBank {
			binary.Write(data, binary.LittleEndian, uint8(0x0A))
		} else {
			binary.Write(data, binary.LittleEndian, uint8(0x0))
		}

		// 2000-2FFF - 8 least significant bits of ROM bank number
		binary.Write(data, binary.LittleEndian, uint16(0x2000))
		binary.Write(data, binary.LittleEndian, uint8(mbc.romBank))

		// 3000-3FFF - 9th bit of ROM bank number
		binary.Write(data, binary.LittleEndian, uint16(0x3000))
		binary.Write(data, binary.LittleEndian, uint8(mbc.romBank>>8))

		// 4000-5FFF - RAM bank number
		binary.Write(data, binary.LittleEndian, uint16(0x4000))
		binary.Write(data, binary.LittleEndian, mbc.ramBank)
	}
}

func (e *Emulator) BessStore(filename string) error {
	data := new(bytes.Buffer)
	romFeatures := e.mem.rom.features

	// RAM
	ram := e.mem.workRam[:]
	sizeRAM := len(ram)
	offsetRAM := 0
	data.Write(ram)

	// VRAM
	vram := e.mem.videoRam[:]
	sizeVRAM := len(vram)
	offsetVRAM := offsetRAM + sizeRAM
	data.Write(vram)

	// MBC RAM
	mbcRam := e.mem.rom.controller.GetRam()
	sizeMBCRAM := len(mbcRam)
	offsetMBCRAM := offsetVRAM + sizeVRAM
	data.Write(mbcRam)

	// OAM
	oam := e.mem.io[:160]
	sizeOAM := len(oam)
	offsetOAM := offsetMBCRAM + sizeMBCRAM
	data.Write(oam)

	// HRAM
	hram := e.mem.io[0x180:0x1ff]
	sizeHRAM := len(hram)
	offsetHRAM := offsetOAM + sizeOAM
	data.Write(hram)

	// Name Block
	data.WriteString("NAME")
	binary.Write(data, binary.LittleEndian, uint32(15))
	data.WriteString("EMULATOR-GO 0.1")

	// Info Block
	data.WriteString("INFO")
	binary.Write(data, binary.LittleEndian, uint32(0x12))
	data.Write(romFeatures.TitleBytes)          // ROM (Title)
	data.Write(romFeatures.GlobalChecksumBytes) // ROM (Global checksum)

	// Core Block
	data.WriteString("CORE")
	binary.Write(data, binary.LittleEndian, uint32(0xD0))
	binary.Write(data, binary.LittleEndian, uint16(1)) // Major BESS version as a 16-bit integer
	binary.Write(data, binary.LittleEndian, uint16(1)) // Major Minor version as a 16-bit integer
	data.WriteString("GDA ")                           // A four-character ASCII model identifier
	binary.Write(data, binary.LittleEndian, e.cpu.PC)
	binary.Write(data, binary.LittleEndian, e.cpu.AF.value)
	binary.Write(data, binary.LittleEndian, e.cpu.BC.value)
	binary.Write(data, binary.LittleEndian, e.cpu.DE.value)
	binary.Write(data, binary.LittleEndian, e.cpu.HL.value)
	binary.Write(data, binary.LittleEndian, e.cpu.SP.value)
	binary.Write(data, binary.LittleEndian, e.IME)
	binary.Write(data, binary.LittleEndian, e.mem.GetIE()) // The value of the IE register
	binary.Write(data, binary.LittleEndian, e.halt)        // Execution state (0 = running; 1 = halted; 2 = stopped)
	binary.Write(data, binary.LittleEndian, uint8(0))      // Reserved, must be 0
	data.Write(e.mem.io[0x100:0x180])                      // Memory-mapped Registers

	binary.Write(data, binary.LittleEndian, uint32(sizeRAM))      // The size of RAM (32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(offsetRAM))    // The offset of RAM from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(sizeVRAM))     // The size of VRAM (32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(offsetVRAM))   // The offset of VRAM from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(sizeMBCRAM))   // The size of MBC RAM (32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(offsetMBCRAM)) // The offset of MBC RAM from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(sizeOAM))      // The size of OAM (=0xA0, 32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(offsetOAM))    // The offset of OAM from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(sizeHRAM))     // The size of HRAM (=0x7F, 32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(offsetHRAM))   // The offset of HRAM from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(0x0))          // The size of background palettes (=0x40 or 0, 32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(0))            // The offset of background palettes from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(0x0))          // The size of object palettes (=0x40 or 0, 32-bit integer)
	binary.Write(data, binary.LittleEndian, uint32(0))            // The offset of object palettes from file start (32-bit integer)

	// XOAM block - Not implemented
	// MBC block - Not implemented
	WriteMBCBlock(data, e.mem.rom.controller)
	// RTC block - Not implemented
	// HUC3 block - Not implemented
	// TPP1 block - Not implemented
	// MBC7 block - Not implemented
	// SGB block - Not implemented

	// End Block
	data.WriteString("END ")
	binary.Write(data, binary.LittleEndian, uint32(0))

	// Footer
	binary.Write(data, binary.LittleEndian, uint32(offsetHRAM+sizeHRAM))
	data.WriteString("BESS")

	err := os.WriteFile(filename, data.Bytes(), 0644)
	if err != nil {
		log.Println("Error creating save state:", err)
		return err
	}

	return nil
}

type BessBlock struct {
	name string
	size uint32
}

type BessInfo struct {
	nameBlockParsed, infoBlockParsed bool

	startOfFile                                                                                                uint32
	emulatorName                                                                                               []byte
	romTitle                                                                                                   [16]byte
	romGlobalChecksum                                                                                          uint16
	reserved                                                                                                   uint8
	majorVersion, minorVersion                                                                                 int16
	modelIdentifier                                                                                            [4]byte
	pc, af, bc, de, hl, sp                                                                                     uint16
	ime, ie, halt                                                                                              uint8
	io                                                                                                         [128]byte
	ramSize, vramSize, mbcRamSize, OAMSize, HRAMSize, BackgroundPalettesSize, ObjectPalettesSize               uint32
	ramOffset, vramOffset, mbcRamOffset, OAMOffset, HRAMOffset, BackgroundPalettesOffset, ObjectPalettesOffset uint32
	extraOAM                                                                                                   [96]byte

	mbcRegisters map[uint16]uint8
}

func (b *BessBlock) PrintHeader() {
	fmt.Println("BLOCK: ", b.name, b.size)
}

func parseBlock(reader io.Reader) (*BessBlock, error) {
	var block BessBlock

	// Block Name
	var blockNameTmp = make([]byte, 4)
	_, err := reader.Read(blockNameTmp)
	if err != nil {
		return nil, err
	}
	block.name = string(blockNameTmp)

	// Block Size
	err = binary.Read(reader, binary.LittleEndian, &block.size)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (bess *BessInfo) parseFooter(bs []byte) error {
	// The ASCII string 'BESS'
	bessString := string(bs[len(bs)-4:])

	if bessString != "BESS" {
		return fmt.Errorf("bess not in footer")
	}

	// Offset to the first BESS Block, from the file's start
	bess.startOfFile = binary.LittleEndian.Uint32(bs[len(bs)-8 : len(bs)-4])

	return nil
}

func (bess *BessInfo) parseCoreBlock(block *BessBlock, reader io.Reader) error {
	if block.size != 0xD0 {
		return fmt.Errorf("incorrect CORE block size: %d", block.size)
	}

	// Major & minor version
	err := binary.Read(reader, binary.LittleEndian, &bess.majorVersion)
	if err != nil {
		return err
	}

	err = binary.Read(reader, binary.LittleEndian, &bess.minorVersion)
	if err != nil {
		return err
	}

	// Model identifier
	_, err = reader.Read(bess.modelIdentifier[:])
	if err != nil {
		return err
	}

	// CPU Registers
	err = binary.Read(reader, binary.LittleEndian, &bess.pc)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.af)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.bc)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.de)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.hl)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.sp)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.ime)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.ie)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.halt)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.reserved)
	if err != nil {
		return err
	}

	// IO Registers
	_, err = reader.Read(bess.io[:])
	if err != nil {
		return err
	}

	// Sizes & Offsets
	err = binary.Read(reader, binary.LittleEndian, &bess.ramSize)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.ramOffset)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.vramSize)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.vramOffset)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.mbcRamSize)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.mbcRamOffset)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.OAMSize)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.OAMOffset)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.HRAMSize)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.HRAMOffset)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.BackgroundPalettesSize)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.BackgroundPalettesOffset)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.ObjectPalettesSize)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &bess.ObjectPalettesOffset)
	if err != nil {
		return err
	}

	// Error checks
	if bess.majorVersion != 1 {
		return fmt.Errorf("major version not supported: %d", bess.majorVersion)
	}

	if bess.modelIdentifier[0] != 'G' {
		return fmt.Errorf("model identifier not supported: %c", bess.modelIdentifier)
	}

	if bess.reserved != 0 {
		return fmt.Errorf("reserved byte with offset 0x17 0x%x != 0", bess.reserved)
	}

	if bess.ime > 1 {
		return fmt.Errorf("incorrect ime: %d", bess.ime)
	}

	if bess.halt > 2 {
		return fmt.Errorf("incorrect execution state: %d", bess.halt)
	}

	if bess.ramSize != 0x2000 {
		return fmt.Errorf("incorrect ram size: %d", bess.ramSize)
	}

	if bess.vramSize != 0x2000 {
		return fmt.Errorf("incorrect vram size: %d", bess.vramSize)
	}

	if bess.OAMSize != 0xA0 {
		return fmt.Errorf("incorrect oam size: %d", bess.OAMSize)
	}

	if bess.HRAMSize != 0x7F {
		return fmt.Errorf("incorrect hram size: %d", bess.HRAMSize)
	}

	if bess.BackgroundPalettesSize != 0x40 && bess.BackgroundPalettesSize != 0x0 {
		return fmt.Errorf("incorrect background palettes size: %d", bess.BackgroundPalettesSize)
	}

	if bess.ObjectPalettesSize != 0x40 && bess.ObjectPalettesSize != 0x0 {
		return fmt.Errorf("incorrect object palettes size: %d", bess.ObjectPalettesSize)
	}

	return nil
}

func (bess *BessInfo) parseNameBlock(block *BessBlock, reader io.Reader) error {
	bess.nameBlockParsed = true
	bess.emulatorName = make([]byte, block.size)
	_, err := reader.Read(bess.emulatorName)
	if err != nil {
		return err
	}

	return nil
}

func (bess *BessInfo) parseInfoBlock(block *BessBlock, reader io.Reader) error {
	bess.nameBlockParsed = true
	if block.size != 18 {
		return fmt.Errorf("incorrect INFO block size: %d", block.size)
	}

	_, err := reader.Read(bess.romTitle[:])
	if err != nil {
		return err
	}

	err = binary.Read(reader, binary.BigEndian, &bess.romGlobalChecksum)
	if err != nil {
		return err
	}

	return nil
}

func (bess *BessInfo) parseMBCBlock(block *BessBlock, reader io.Reader) error {
	if block.size%3 != 0 {
		return fmt.Errorf("incorrect block size %d. Must be divisible by 3 in MBC block", block.size)
	}

	bess.mbcRegisters = make(map[uint16]uint8)
	numRegisters := block.size / 3
	for i := uint32(0); i < numRegisters; i++ {
		var addr uint16
		var value uint8

		err := binary.Read(reader, binary.LittleEndian, &addr)
		if err != nil {
			return err
		}
		err = binary.Read(reader, binary.LittleEndian, &value)
		if err != nil {
			return err
		}

		bess.mbcRegisters[addr] = value
	}

	return nil
}

func (bess *BessInfo) parseXOAMBlock(block *BessBlock, reader io.Reader) error {
	if block.size != 0x60 {
		return fmt.Errorf("incorrect block size %d. Must be 0x60 in MBC block", block.size)
	}

	_, err := reader.Read(bess.extraOAM[:])
	if err != nil {
		return err
	}

	return nil
}

func (e *Emulator) BessLoad(filename string) error {
	romFeatures := e.mem.rom.features
	var bess BessInfo
	var blocksParsed = make(map[string]bool)

	bs, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Error creating save state:", err)
		return err
	}

	// Footer
	err = bess.parseFooter(bs)
	if err != nil {
		return err
	}

	buffer := bytes.NewReader(bs[bess.startOfFile : len(bs)-8])

	// Name Block & Info Blocks
coreBlockLoop:
	for {
		block, err := parseBlock(buffer)
		if err != nil {
			return err
		}
		block.PrintHeader()

		// Check for duplicate blocks
		if _, ok := blocksParsed[block.name]; ok {
			return fmt.Errorf("duplicate block '%s'", block.name)
		}
		blocksParsed[block.name] = true

		switch block.name {
		case "NAME":
			if bess.infoBlockParsed {
				return fmt.Errorf("incorrect block position. 'NAME' block has to be before 'INFO' block")
			}

			err = bess.parseNameBlock(block, buffer)
			if err != nil {
				return err
			}
		case "INFO":
			err = bess.parseInfoBlock(block, buffer)
			if err != nil {
				return err
			}
		case "CORE":
			err = bess.parseCoreBlock(block, buffer)
			if err != nil {
				return err
			}
			break coreBlockLoop
		default:
			return fmt.Errorf("incorrect block %s", block.name)
		}
	}

blockLoop:
	for {
		block, err := parseBlock(buffer)
		if err != nil {
			return err
		}

		// Check for duplicate blocks
		if _, ok := blocksParsed[block.name]; ok {
			return fmt.Errorf("duplicate block '%s'", block.name)
		}
		blocksParsed[block.name] = true
		block.PrintHeader()

		switch block.name {
		case "XOAM":
			err = bess.parseXOAMBlock(block, buffer)
			if err != nil {
				return err
			}
		case "MBC ":
			err = bess.parseMBCBlock(block, buffer)
			if err != nil {
				return err
			}
		case "RTC ":
			content := make([]byte, block.size)
			_, err = buffer.Read(content)
		case "HUC3":
			content := make([]byte, block.size)
			_, err = buffer.Read(content)
		case "TPP1":
			content := make([]byte, block.size)
			_, err = buffer.Read(content)
		case "MBC7":
			content := make([]byte, block.size)
			_, err = buffer.Read(content)
		case "SGB ":
			content := make([]byte, block.size)
			_, err = buffer.Read(content)
		case "END ":
			if block.size != 0 {
				return fmt.Errorf("bess end block size (%d) != 0", block.size)
			}
			break blockLoop
		case "NAME":
			return fmt.Errorf("incorrect block position. 'NAME' block has to be before 'CORE' block")
		case "INFO":
			return fmt.Errorf("incorrect block position. 'INFO' block has to be before 'CORE' block")
		default:
			return fmt.Errorf("incorrect block '%s'", block.name)
		}
	}

	// Extra checks
	if res := bytes.Compare(bess.romTitle[:], romFeatures.TitleBytes); res != 0 {
		return fmt.Errorf("incorrect rom title: %v != %v", bess.romTitle, []byte(romFeatures.Title))
	}

	if bess.romGlobalChecksum != romFeatures.GlobalChecksum {
		return fmt.Errorf("incorrect rom checksum: %d != %d", bess.romGlobalChecksum, romFeatures.GlobalChecksum)
	}

	// Set Register values
	e.cpu.PC = bess.pc
	e.cpu.SetAF(bess.af)
	e.cpu.SetBC(bess.bc)
	e.cpu.SetDE(bess.de)
	e.cpu.SetHL(bess.hl)
	e.cpu.SetSP(bess.sp)
	e.mem.SetIE(bess.ie)
	e.IME = bess.ime
	e.halt = bess.halt

	// Copy I/O Registers
	copy(e.mem.io[0x100:], bess.io[:])

	div := e.mem.GetDIV()
	e.timer.SetInternalTimer(uint16(div) << 8)

	// Copy Extra OAM
	copy(e.mem.io[0xA0:], bess.extraOAM[:])

	// Write MBC Registers
	for addr, value := range bess.mbcRegisters {
		e.mem.write8(addr, value)
	}

	// RAM
	copy(e.mem.workRam[:], bs[bess.ramOffset:bess.ramOffset+bess.ramSize])

	// VRAM
	copy(e.mem.videoRam[:], bs[bess.vramOffset:bess.vramOffset+bess.vramSize])

	// MBC RAM
	copy(e.mem.rom.controller.GetRam(), bs[bess.mbcRamOffset:bess.mbcRamOffset+bess.mbcRamSize])

	// OAM
	copy(e.mem.io[:], bs[bess.OAMOffset:bess.OAMOffset+bess.OAMSize])

	// HRAM
	copy(e.mem.io[0x180:], bs[bess.HRAMOffset:bess.HRAMOffset+bess.HRAMSize])

	// background palettes

	// object palettes

	return nil
}
