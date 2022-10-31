package emulator

import (
	"emulator-go/emulator-old/gb/utils"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"time"
)

const WIDTH = 160
const HEIGHT = 144

type Emulator struct {
	cycles, prevCycles uint16
	totalCycles        uint64

	IME  uint8
	halt uint8

	workRam  [0x4000]uint8
	videoRam [0x2000]uint8

	romFilename     string
	bootRomFilename string

	io                [0x200]uint8
	extrambank        *[0x8000]uint8
	ppuDot            int
	rom0              []byte
	bootRom           []byte
	extrambankPointer uint32
	rom1Pointer       uint32
	keyboardState     []uint8
	frameBuffer       [WIDTH * HEIGHT]int32
	lcdcControl       LCDControl

	palette []int32

	cpu CPU

	window                    *sdl.Window
	renderer                  *sdl.Renderer
	texture                   *sdl.Texture
	font                      *ttf.Font
	frames                    uint64
	framesPerSecond           uint32
	framesCurrentSecond       uint32
	deltaTime                 uint64
	millisecondsPreviousFrame uint64
	consoleMessage            string
	showMessage               bool
	consoleMessageDuration    time.Duration
	consoleMessageStart       time.Time

	numInstructions      uint64
	vsyncEnabled         bool
	showFPS              bool
	stop                 bool
	pause                bool
	reset                bool
	bootRomEnabled       bool
	romHeader            RomHeader
	memoryBankController int
}

func NewEmulator(romFilename, saveFilename, bootRomFilename string) (*Emulator, error) {
	emulator := Emulator{
		ppuDot:          32,
		rom1Pointer:     32768,
		palette:         []int32{-1, -23197, -65536, -1 << 24, -1, -8092417, -12961132, -1 << 24},
		vsyncEnabled:    true,
		showFPS:         false,
		showMessage:     false,
		romFilename:     romFilename,
		bootRomFilename: bootRomFilename,
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

	err := emulator.loadRom(romFilename)
	if err != nil {
		return nil, err
	}

	err = emulator.initializeSDL(utils.ToCamel(emulator.romHeader.Title), 4)
	if err != nil {
		return nil, err
	}

	err = emulator.initializeSaveFile(saveFilename)
	if err != nil {
		return nil, err
	}

	return &emulator, nil
}

func (e *Emulator) Destroy() {
	e.texture.Destroy()
	e.renderer.Destroy()
	e.window.Destroy()
	sdl.Quit()
}

func (e *Emulator) Run() {
	e.stop = false
	for {
		if e.reset {
			e.Reset()
		}

		e.prevCycles = e.cycles
		if (e.IME & e.GetIF() & e.io[511]) != 0 {
			e.SetIF(0)
			e.halt = 0
			e.IME = 0
			e.tick()
			e.tick()
			e.push(e.cpu.PC)
			e.cpu.PC = 64

			e.loadGamesharkCodes()

		} else if e.halt != 0 {
			e.tick()
		} else {
			e.CPURun()
		}

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

	err := e.loadRom(e.romFilename)
	if err != nil {
		return err
	}

	return nil
}
