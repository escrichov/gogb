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

	workRam  [0x4000]uint8
	videoRam [0x2000]uint8

	romFilename     string
	bootRomFilename string

	io      [0x200]uint8
	ppuDot  int
	bootRom []byte

	// Rom
	rom *Rom

	lcdcControl LCDControl
	lcdStatus   LCDStatus

	// Timers
	internalTimer                  uint16
	timaUpdateWithTMADelayedCycles uint64

	// Serial
	serialTransferedBits uint8

	palette []int32

	cpu CPU

	// SDL Window
	window *Window

	numInstructions uint64
	stop            bool
	pause           bool
	reset           bool
	bootRomEnabled  bool
}

func NewEmulator(romFilename, bootRomFilename, fontFilename string, showWindow bool) (*Emulator, error) {
	var err error

	emulator := Emulator{
		ppuDot:          32,
		palette:         []int32{-1, -23197, -65536, -1 << 24, -1, -8092417, -12961132, -1 << 24},
		romFilename:     romFilename,
		bootRomFilename: bootRomFilename,
		internalTimer:   8,
	}

	if bootRomFilename == "" {
		emulator.initializeBootRomValues()
		emulator.bootRomEnabled = false
	} else {
		err := emulator.loadBootRom(bootRomFilename)
		if err != nil {
			return nil, err
		}
		emulator.bootRomEnabled = true
	}

	emulator.rom, err = newRomFromFile(romFilename)
	if err != nil {
		return nil, err
	}

	emulator.window, err = newWindow(
		ToCamel(emulator.rom.features.Title),
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

		e.incrementTimers()
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

		e.incrementTimers()
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

	if e.bootRomEnabled {
		err := e.loadBootRom(e.bootRomFilename)
		if err != nil {
			return err
		}
	} else {
		e.initializeBootRomValues()
	}

	// Framebuffer set to black
	e.window.SetFramebufferColor(0)

	e.rom, err = newRomFromFile(e.romFilename)
	if err != nil {
		return err
	}

	return nil
}
