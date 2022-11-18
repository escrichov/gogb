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
}

func NewEmulator(romFilename, bootRomFilename, fontFilename string, showWindow bool) (*Emulator, error) {
	var err error

	emulator := Emulator{
		romFilename:     romFilename,
		bootRomFilename: bootRomFilename,
		ppu: PPU{
			ppuDot:     32,
			paletteBGP: []uint32{0xFFFFFFFF, 0xFFFFA563, 0xFFFF0000, 0xFF000000},
			paletteOB0: []uint32{0xFFFFFFFF, 0xFF8484FF, 0xFF3A3A94, 0xFF000000},
			paletteOB1: []uint32{0xFFFFFFFF, 0xFF8484FF, 0xFF3A3A94, 0xFF000000},
		},
		timer: Timer{internalTimer: 8},
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
	)

	emulator.PrintCartridge()

	return &emulator, nil
}

func (e *Emulator) Destroy() {
	e.window.Destroy()
}

func (e *Emulator) RunTest(numCycles uint64) {
	e.stop = false
	for {
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

		renderFrame := e.PPURun()
		if renderFrame {
			e.window.renderFrame()
		}

		// Paused state
		for e.pause {
			e.window.renderFrame()
		}

		if e.stop {
			break
		}

		if numCycles != 0 && e.cycles >= numCycles {
			break
		}
	}
}

func (e *Emulator) Run() {
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

		renderFrame := e.PPURun()
		if renderFrame {
			e.window.renderFrame()
			e.manageKeyboardEvents()
		}

		// Paused state
		for e.pause {
			e.window.renderFrame()
			e.manageKeyboardEvents()
			time.Sleep(time.Millisecond * 1000 / 60.0)
		}

		if e.stop {
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
