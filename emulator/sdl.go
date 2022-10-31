package emulator

import (
	"bytes"
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
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

func pngToJpeg(inputPNG []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(inputPNG))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func SaveJPGFromSurface(surface *sdl.Surface, filename string) error {
	var imagePng = make([]byte, 10000)
	rw, err := sdl.RWFromMem(imagePng)
	if err != nil {
		return err
	}
	defer rw.Close()

	err = img.SavePNGRW(surface, rw, 0)
	if err != nil {
		return err
	}

	nTell, err := rw.Tell()
	if err != nil {
		return err
	}

	imageJpeg, err := pngToJpeg(imagePng[:nTell])
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, imageJpeg, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (e *Emulator) TakeSnapshot(filename string) error {
	w, h, err := e.renderer.GetOutputSize()
	if err != nil {
		return err
	}

	surface, err := sdl.CreateRGBSurface(0, w, h, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
	if err != nil {
		return err
	}

	err = e.renderer.ReadPixels(nil, sdl.PIXELFORMAT_ARGB8888, surface.Data(), int(surface.Pitch))
	if err != nil {
		return err
	}

	filenameExtension := filepath.Ext(filename)
	if filenameExtension == ".bmp" {
		err = surface.SaveBMP(filename)
		if err != nil {
			return err
		}
	} else if filenameExtension == ".png" {
		err = img.SavePNG(surface, filename)
		if err != nil {
			return err
		}
	} else if filenameExtension == ".jpeg" || filenameExtension == ".jpg" {
		err = SaveJPGFromSurface(surface, filename)
		if err != nil {
			return err
		}
	} else {
		fmt.Errorf("file extension not recognized: %s", filenameExtension)
	}

	surface.Free()

	return nil
}
