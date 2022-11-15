package emulator

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func (e *Emulator) BessStore(filename string) error {
	data := new(bytes.Buffer)

	// Name Block
	data.WriteString("NAME")
	binary.Write(data, binary.LittleEndian, int32(15))
	data.WriteString("EMULATOR-GO 0.1")

	// Info Block
	data.WriteString("INFO")
	binary.Write(data, binary.LittleEndian, int32(0x12))
	data.Write(e.rom.features.TitleBytes)          // ROM (Title)
	data.Write(e.rom.features.GlobalChecksumBytes) // ROM (Global checksum)

	// Core Block
	data.WriteString("CORE")
	binary.Write(data, binary.LittleEndian, int32(0xD0))
	binary.Write(data, binary.LittleEndian, int16(1)) // Major BESS version as a 16-bit integer
	binary.Write(data, binary.LittleEndian, int16(1)) // Major Minor version as a 16-bit integer
	data.WriteString("GDA ")                          // A four-character ASCII model identifier
	binary.Write(data, binary.LittleEndian, e.cpu.PC)
	binary.Write(data, binary.LittleEndian, e.cpu.AF.value)
	binary.Write(data, binary.LittleEndian, e.cpu.BC.value)
	binary.Write(data, binary.LittleEndian, e.cpu.DE.value)
	binary.Write(data, binary.LittleEndian, e.cpu.HL.value)
	binary.Write(data, binary.LittleEndian, e.cpu.SP.value)
	binary.Write(data, binary.LittleEndian, e.IME)
	binary.Write(data, binary.LittleEndian, e.GetIF()) // The value of the IE register
	binary.Write(data, binary.LittleEndian, e.halt)    // Execution state (0 = running; 1 = halted; 2 = stopped)
	binary.Write(data, binary.LittleEndian, uint8(0))  // Reserved, must be 0
	data.Write(e.io[0x100:0x180])                      // Memory-mapped Registers

	binary.Write(data, binary.LittleEndian, int32(0x4000)) // The size of RAM (32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of RAM from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0x2000)) // The size of VRAM (32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of VRAM from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0))      // The size of MBC RAM (32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of MBC RAM from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0xA0))   // The size of OAM (=0xA0, 32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of OAM from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0x7F))   // The size of HRAM (=0x7F, 32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of HRAM from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0x40))   // The size of background palettes (=0x40 or 0, 32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of background palettes from file start (32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0x40))   // The size of object palettes (=0x40 or 0, 32-bit integer)
	binary.Write(data, binary.LittleEndian, int32(0))      // The offset of object palettes from file start (32-bit integer)

	// XOAM block - Not implemented
	// MBC block - Not implemented
	// RTC block - Not implemented
	// HUC3 block - Not implemented
	// TPP1 block - Not implemented
	// MBC7 block - Not implemented
	// SGB block - Not implemented

	// End Block
	data.WriteString("END ")
	binary.Write(data, binary.LittleEndian, int32(0))

	// Footer
	binary.Write(data, binary.LittleEndian, int32(0))
	data.WriteString("BESS")

	err := os.WriteFile(filename, data.Bytes(), 0644)
	if err != nil {
		log.Println("Error creating save state:", err)
		return err
	}

	return nil
}

func (e *Emulator) BessLoad(filename string) error {
	bs, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Error creating save state:", err)
		return err
	}

	// Footer
	bess := string(bs[len(bs)-4:])
	startOfFile := binary.LittleEndian.Uint32(bs[len(bs)-8 : len(bs)-4])

	if bess != "BESS" {
		return fmt.Errorf("bess not in footer")
	}

	buffer := bytes.NewReader(bs[startOfFile : len(bs)-8])

	var blockNameTmp []byte = make([]byte, 4)

	// Name Block
	_, err = buffer.Read(blockNameTmp)
	if err != nil {
		return err
	}
	blockName := string(blockNameTmp)

	var blockSize int32
	err = binary.Read(buffer, binary.LittleEndian, &blockSize)
	if err != nil {
		return err
	}

	var blockContent []byte = make([]byte, 15)
	_, err = buffer.Read(blockContent)
	if err != nil {
		return err
	}
	fmt.Println(blockName, blockSize, string(blockContent))

	// Info  Block
	_, err = buffer.Read(blockNameTmp)
	if err != nil {
		return err
	}
	blockName = string(blockNameTmp)

	err = binary.Read(buffer, binary.LittleEndian, &blockSize)
	if err != nil {
		return err
	}

	blockContent = make([]byte, blockSize)
	_, err = buffer.Read(blockContent)
	if err != nil {
		return err
	}
	romTitle := blockContent[:16]
	romGlobalChecksum := binary.BigEndian.Uint16(blockContent[16:])
	fmt.Println(blockName, blockSize, string(romTitle), romGlobalChecksum)

	// Core Block
	_, err = buffer.Read(blockNameTmp)
	if err != nil {
		return err
	}
	blockName = string(blockNameTmp)

	err = binary.Read(buffer, binary.LittleEndian, &blockSize)
	if err != nil {
		return err
	}
	fmt.Println(blockName, blockSize)

	var majorVersion int16
	err = binary.Read(buffer, binary.LittleEndian, &majorVersion)
	if err != nil {
		return err
	}

	var minorVersion int16
	err = binary.Read(buffer, binary.LittleEndian, &minorVersion)
	if err != nil {
		return err
	}
	fmt.Printf("BESS Version: %d.%d\n", majorVersion, minorVersion)

	var modelIdentifier []byte = make([]byte, 4)
	_, err = buffer.Read(modelIdentifier)
	if err != nil {
		return err
	}
	fmt.Println("Model Identifier:", string(modelIdentifier))

	var pc, af, bc, de, hl, sp uint16
	var ime, ie, halt uint8

	err = binary.Read(buffer, binary.LittleEndian, &pc)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &af)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &bc)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &de)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &hl)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &sp)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &ime)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &ie)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &halt)
	if err != nil {
		return err
	}
	reserved, err := buffer.ReadByte()
	if err != nil {
		return err
	}

	fmt.Printf("PC: 0x%x\n", pc)
	fmt.Printf("AF: 0x%x\n", af)
	fmt.Printf("BC: 0x%x\n", bc)
	fmt.Printf("DE: 0x%x\n", de)
	fmt.Printf("HL: 0x%x\n", hl)
	fmt.Printf("SP: 0x%x\n", sp)
	fmt.Printf("IME: %d\n", ime)
	fmt.Printf("IE: %d\n", ie)
	fmt.Printf("HALT: %d\n", halt)

	var io [128]byte
	_, err = buffer.Read(io[:])
	if err != nil {
		return err
	}

	var ramSize, vramSize, mbcRamSize, OAMSize, HRAMSize, BackgroundPalettesSize, ObjectPalettesSize int32
	var ramOffset, vramOffset, mbcRamOffset, OAMOffset, HRAMOffset, BackgroundPalettesOffset, ObjectPalettesOffset int32
	err = binary.Read(buffer, binary.LittleEndian, &ramSize)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &ramOffset)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &vramSize)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &vramOffset)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &mbcRamSize)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &mbcRamOffset)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &OAMSize)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &OAMOffset)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &HRAMSize)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &HRAMOffset)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &BackgroundPalettesSize)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &BackgroundPalettesOffset)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &ObjectPalettesSize)
	if err != nil {
		return err
	}
	err = binary.Read(buffer, binary.LittleEndian, &ObjectPalettesOffset)
	if err != nil {
		return err
	}

	for {
		_, err = buffer.Read(blockNameTmp)
		if err != nil {
			return err
		}
		blockName = string(blockNameTmp)

		err = binary.Read(buffer, binary.LittleEndian, &blockSize)
		if err != nil {
			return err
		}

		fmt.Println([]byte(blockName), blockName, blockSize)

		if blockName == "END " {
			if blockSize != 0 {
				return fmt.Errorf("bess end block size (%d) != 0", blockSize)
			}
			fmt.Println(blockName, blockSize)
			break
		} else {
			blockContent = make([]byte, blockSize)
			_, err = buffer.Read(blockContent)
			if err != nil {
				return err
			}
		}
	}

	// Extra checks
	if res := bytes.Compare(romTitle, e.rom.features.TitleBytes); res != 0 {
		return fmt.Errorf("incorrect rom title: %v != %v", []byte(romTitle), []byte(e.rom.features.Title))
	}

	if romGlobalChecksum != e.rom.features.GlobalChecksum {
		return fmt.Errorf("incorrect rom checksum: %d != %d", romGlobalChecksum, e.rom.features.GlobalChecksum)
	}

	if majorVersion != 1 {
		return fmt.Errorf("major version not supported: %d", majorVersion)
	}

	if modelIdentifier[0] != 'G' {
		return fmt.Errorf("model identifier not supported: %c", modelIdentifier)
	}

	if reserved != 0 {
		return fmt.Errorf("reserved byte with offset 0x17 0x%x != 0", reserved)
	}

	if ime > 1 {
		return fmt.Errorf("incorrect ime: %d", ime)
	}

	if halt > 2 {
		return fmt.Errorf("incorrect execution state: %d", halt)
	}

	if ramSize != 0x2000 {
		return fmt.Errorf("incorrect ram size: %d", ramSize)
	}

	if vramSize != 0x2000 {
		return fmt.Errorf("incorrect vram size: %d", vramSize)
	}

	if OAMSize != 0xA0 {
		return fmt.Errorf("incorrect oam size: %d", OAMSize)
	}

	if HRAMSize != 0x7F {
		return fmt.Errorf("incorrect hram size: %d", HRAMSize)
	}

	if BackgroundPalettesSize != 0x40 {
		return fmt.Errorf("incorrect background palettes size: %d", BackgroundPalettesSize)
	}

	if ObjectPalettesSize != 0x40 {
		return fmt.Errorf("incorrect object palettes size: %d", ObjectPalettesSize)
	}

	// Set values
	e.cpu.PC = pc
	e.cpu.SetAF(af)
	e.cpu.SetBC(bc)
	e.cpu.SetDE(de)
	e.cpu.SetHL(hl)
	e.cpu.SetSP(sp)
	e.SetIF(ie)
	e.IME = ime
	e.halt = halt

	return nil
}
