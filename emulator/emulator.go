package emulator

import (
	"time"
)

const WIDTH = 160
const HEIGHT = 144

type Emulator struct {
	cycles, prevCycles uint64

	IME                                    uint8
	delayedActivateIMEatInstruction        uint64
	halt                                   uint8
	isInterruptPendingInFirstHaltExecution bool
	isIMEDelayedInFirstHaltExecution       bool
	isHaltBugActive                        bool
	isHaltBugEIActive                      bool

	romFilename     string
	bootRomFilename string

	cpu    CPU
	mem    Memory
	ppu    PPU
	timer  Timer
	window *Window

	GameSharkCodes []*GameSharkCode
	GameGenieCodes []*GameGenieCode

	// Serial
	serialTransferedBits uint8

	numInstructions uint64
	stop            bool
	pause           bool
	reset           bool
	showWindow      bool
}

func NewEmulator(romFilename, bootRomFilename, fontFilename string, showWindow, vsyncAtStartup bool) (*Emulator, error) {
	var err error

	emulator := Emulator{
		romFilename:     romFilename,
		bootRomFilename: bootRomFilename,
		ppu: PPU{
			ppuDot:     32,
			paletteBGP: []uint32{0xFFFFFFFF, 0xFFFFA563, 0xFFFF0000, 0xFF000000},
			paletteOB0: []uint32{0xFFFFFFFF, 0xFF8484FF, 0xFF3A3A94, 0xFF000000},
			paletteOB1: []uint32{0xFFFFFFFF, 0xFFFFA563, 0xFFFF0000, 0xFF000000},
		},
		timer:      Timer{internalTimer: 8},
		showWindow: showWindow,
	}
	emulator.mem.emulator = &emulator
	emulator.timer.emulator = &emulator

	if bootRomFilename == "" {
		emulator.initializeBootRomValues()
	} else {
		err := emulator.mem.loadBootRom(bootRomFilename)
		if err != nil {
			return nil, err
		}
	}

	emulator.mem.rom, err = newRomFromFile(romFilename)
	if err != nil {
		return nil, err
	}

	emulator.window, err = newWindow(
		ToCamel(emulator.mem.rom.features.Title),
		fontFilename,
		4.0,
		showWindow,
		vsyncAtStartup,
	)

	emulator.PrintCartridge()

	return &emulator, nil
}

func (e *Emulator) Destroy() {
	e.window.Destroy()
}

func (e *Emulator) Run(numCycles uint64) {
	e.stop = false
	for {
		if e.reset {
			e.Reset()
		}

		e.prevCycles = e.cycles
		if e.hastToManageInterrupts() {
			e.manageInterrupts()
		} else if e.halt != 0 {
			e.HaltRun()
		} else {
			e.CPURun()
		}

		e.timer.incrementTimers(uint16(e.cycles - e.prevCycles))
		e.serialTransfer()

		e.PPURun()

		// Paused state
		for e.pause {
			e.window.renderFrame()
			if e.showWindow {
				e.manageKeyboardEvents()
			}
			time.Sleep(time.Millisecond * 1000 / 60.0)
		}

		if e.stop {
			break
		}

		if numCycles != 0 && e.cycles >= numCycles {
			break
		}
	}
}

func (e *Emulator) Reset() error {
	var err error
	e.reset = false

	if e.mem.bootRomEnabled {
		err := e.mem.loadBootRom(e.bootRomFilename)
		if err != nil {
			return err
		}
	} else {
		e.initializeBootRomValues()
	}

	// Framebuffer set to black
	e.window.SetFramebufferColor(0)

	e.mem.rom, err = newRomFromFile(e.romFilename)
	if err != nil {
		return err
	}

	return nil
}
