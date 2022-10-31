package emulator_old

import (
	"emulator-go/emulator/gb"
)

type Emulator struct {
	isRunning        bool
	MainWindow       Window
	BackgroundWindow Window
	GB               gb.GB
}

func (e *Emulator) Initialize(bootRoom []byte, gameRom []byte) error {
	e.isRunning = true
	err := e.GB.Initialize(bootRoom, gameRom)
	if err != nil {
		return err
	}

	err = SDLInit()
	if err != nil {
		return err
	}

	err = e.MainWindow.Initialize("Gameboy", 160, 144, 60, false, false)
	if err != nil {
		e.MainWindow.Destroy()
		return err
	}

	err = e.BackgroundWindow.Initialize("Background", 256, 256, 60, false, false)
	if err != nil {
		e.BackgroundWindow.Destroy()
		return err
	}

	return nil
}

func (e *Emulator) Run() {

	for e.isRunning {
		e.GB.Run()
		e.MainWindow.SetColorBuffer(e.GB.PPU.Screen)
		e.BackgroundWindow.SetColorBuffer(e.GB.PPU.BackgroundScreen)
		e.ProcessEvents()
		e.MainWindow.Update()
		e.MainWindow.Render()

		e.BackgroundWindow.Update()
		e.BackgroundWindow.Render()
	}
}

func (e *Emulator) Destroy() {
	e.MainWindow.Destroy()
	e.BackgroundWindow.Destroy()
}
