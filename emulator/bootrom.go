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
	e.SetLCDC(145)
	e.SetLY(0)
	e.SetDIV(44032)
	e.cpu.SetSP(65534)
	e.cpu.SetA(1)
	e.cpu.SetF(176)
	e.cpu.SetB(19)
	e.cpu.SetC(0)
	e.cpu.SetD(0)
	e.cpu.SetE(216)
	e.cpu.SetH(1)
	e.cpu.SetL(77)
}
