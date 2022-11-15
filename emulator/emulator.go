package emulator

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

	keyboardState []uint8
	frameBuffer   [WIDTH * HEIGHT]int32
	lcdcControl   LCDControl
	lcdStatus     LCDStatus

	// Timers
	internalTimer                  uint16
	timaUpdateWithTMADelayedCycles uint64

	// Serial
	serialTransferedBits uint8

	palette []int32

	cpu CPU

	window                    *sdl.Window
	surface                   *sdl.Surface
	renderer                  *sdl.Renderer
	texture                   *sdl.Texture
	font                      *ttf.Font
	showWindow                bool
	frames                    uint64
	framesPerSecond           uint32
	framesCurrentSecond       uint32
	deltaTime                 uint64
	millisecondsPreviousFrame uint64
	consoleMessage            string
	showMessage               bool
	consoleMessageDuration    time.Duration
	consoleMessageStart       time.Time

	numInstructions uint64
	vsyncEnabled    bool
	showFPS         bool
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
		vsyncEnabled:    true,
		showFPS:         false,
		showMessage:     false,
		romFilename:     romFilename,
		bootRomFilename: bootRomFilename,
		showWindow:      showWindow,
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

	// Framebuffer set to black
	for i, _ := range emulator.frameBuffer {
		emulator.frameBuffer[i] = 0
	}

	emulator.rom, err = newRomFromFile(romFilename)
	if err != nil {
		return nil, err
	}

	err = emulator.initializeSDL(ToCamel(emulator.rom.features.Title), fontFilename, 4)
	if err != nil {
		return nil, err
	}

	emulator.PrintCartridge()

	return &emulator, nil
}

func (e *Emulator) Destroy() {
	e.texture.Destroy()
	e.renderer.Destroy()
	e.window.Destroy()
	sdl.Quit()
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
			e.renderFrame()
		}

		// Paused state
		for e.pause {
			e.renderFrame()
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
			e.renderFrame()
			e.manageKeyboardEvents()
		}

		// Paused state
		for e.pause {
			e.renderFrame()
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
	for i, _ := range e.frameBuffer {
		e.frameBuffer[i] = 0
	}

	e.rom, err = newRomFromFile(e.romFilename)
	if err != nil {
		return err
	}

	return nil
}
