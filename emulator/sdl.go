package emulator

import (
	"github.com/veandco/go-sdl2/sdl"
	"unsafe"
)

func (e *Emulator) initializeSDL(windowName string, windowScale float64) error {
	var err error

	// SDL Initialization
	var subsystemMask uint32 = sdl.INIT_VIDEO | sdl.INIT_AUDIO
	if sdl.WasInit(subsystemMask) != subsystemMask {
		if err := sdl.Init(subsystemMask); err != nil {
			return err
		}
	}

	e.window, err = sdl.CreateWindow(
		windowName,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(WIDTH*windowScale),
		int32(HEIGHT*windowScale),
		sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
	if err != nil {
		return err
	}

	e.renderer, err = sdl.CreateRenderer(e.window, -1, sdl.RENDERER_PRESENTVSYNC|sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}

	// Creating a SDL texture that is used to display the color buffer
	e.texture, err = e.renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		WIDTH,
		HEIGHT,
	)
	if err != nil {
		return err
	}

	// Point to Keyboard State
	e.keyboardState = sdl.GetKeyboardState()

	// Vsync
	e.renderer.RenderSetVSync(e.vsyncEnabled)

	return nil
}

func (e *Emulator) renderFrame() {
	e.renderer.Clear()
	buf := unsafe.Pointer(&e.frameBuffer[0])
	framebufferBytes := unsafe.Slice((*byte)(buf), WIDTH*HEIGHT)
	e.texture.Update(nil, framebufferBytes, WIDTH*4)
	e.renderer.Copy(e.texture, nil, nil)
	e.renderer.Present()
}
