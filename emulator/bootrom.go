package emulator

import (
	"log"
	"os"
)

func (e *Emulator) loadBootRom(fileName string) error {
	var err error

	e.bootRom, err = os.ReadFile(fileName)
	if err != nil {
		log.Println("BootRom file not found:", err)
		return err
	}

	return nil
}

func (e *Emulator) initializeBootRomValues() {
	e.cpu.PC = 256
	e.cpu.SetSP(0xfffe)
	e.cpu.SetA(0x01)
	e.cpu.SetF(0xb0)
	e.cpu.SetB(0x00)
	e.cpu.SetC(0x13)
	e.cpu.SetD(0x00)
	e.cpu.SetE(0xd8)
	e.cpu.SetH(0x01)
	e.cpu.SetL(0x4d)

	e.IME = 0
	e.SetLCDC(145)
	e.SetLY(0)
	e.SetDIV(0xab)
	e.SetInternalTimer(0xabcc) // Or 0xabc4 if not initialized in 8 cycles at startup

	//e.cpu.PC = 256
	//e.SetLCDC(145)
	//e.SetLY(0)
	//e.SetDIV(44032)
	//e.cpu.SetSP(0xfffe)
	//e.cpu.SetA(0x11)
	//e.cpu.SetF(0x80)
	//e.cpu.SetB(0x00)
	//e.cpu.SetC(0x00)
	//e.cpu.SetD(0xff)
	//e.cpu.SetE(0x56)
	//e.cpu.SetH(0x00)
	//e.cpu.SetL(0x0d)
}
