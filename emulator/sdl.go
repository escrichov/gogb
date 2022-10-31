package emulator

import (
	"bytes"
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"time"
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

	// Initialize SDL TTF
	if err := ttf.Init(); err != nil {
		return err
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

	// Font
	e.font, err = ttf.OpenFont("assets/fonts/arial.ttf", 25)
	if err != nil {
		return err
	}

	// Vsync
	e.renderer.RenderSetVSync(e.vsyncEnabled)

	return nil
}

func (e *Emulator) calculateFPS() {
	e.frames++
	e.framesCurrentSecond++

	ticks := sdl.GetTicks64()
	e.deltaTime += ticks - e.millisecondsPreviousFrame
	e.millisecondsPreviousFrame = ticks
	if e.deltaTime >= 1000 {
		e.framesPerSecond = e.framesCurrentSecond
		e.framesCurrentSecond = 0
		e.deltaTime = 0
	}
}

func (e *Emulator) renderFrame() {
	e.calculateFPS()

	e.renderer.Clear()

	// Render FrameBuffer
	buf := unsafe.Pointer(&e.frameBuffer[0])
	framebufferBytes := unsafe.Slice((*byte)(buf), WIDTH*HEIGHT)
	e.texture.Update(nil, framebufferBytes, WIDTH*4)
	e.renderer.Copy(e.texture, nil, nil)

	if e.showFPS {
		e.drawFPS()
	}

	if e.showMessage {
		e.drawMessage()

		now := time.Now()
		if now.Sub(e.consoleMessageStart) > e.consoleMessageDuration {
			e.showMessage = false
		}
	}

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

func (e *Emulator) drawFPS() error {
	color := sdl.Color{R: 0, G: 255, B: 0, A: 255}
	framesString := fmt.Sprintf("FPS %d", e.framesPerSecond)
	return e.drawText(WIDTH*4-100, 0, framesString, color)
}

func (e *Emulator) drawMessage() error {
	color := sdl.Color{R: 0, G: 0, B: 255, A: 255}
	return e.drawText(10, HEIGHT*4-35, e.consoleMessage, color)
}

func (e *Emulator) SetMessage(message string, duration time.Duration) {
	e.showMessage = true
	e.consoleMessage = message
	e.consoleMessageDuration = duration
	e.consoleMessageStart = time.Now()
}

func (e *Emulator) drawText(x, y int32, text string, color sdl.Color) error {

	// as TTF_RenderText_Solid could only be used on
	// SDL_Surface then you have to create the surface first
	surface, err := e.font.RenderUTF8Solid(text, color)
	if err != nil {
		return err
	}
	defer surface.Free()

	// now you can convert it into a texture
	texture, err := e.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	textWidth, textHeight, err := e.font.SizeUTF8(text)

	messageRectangle := sdl.Rect{
		X: x,                 // Controls the rect's x coordinate
		Y: y,                 // controls the rect's y coordinate
		W: int32(textWidth),  // controls the width of the rect
		H: int32(textHeight), // controls the height of the rect
	}

	// (0,0) is on the top left of the window/screen,
	// think a rect as the text's box,
	// that way it would be very simple to understand

	// Now since it's a texture, you have to put RenderCopy
	// in your game loop area, the area where the whole code executes

	// you put the renderer's name first, the Message,
	// the crop size (you can ignore this if you don't want
	// to dabble with cropping), and the rect which is the size
	// and coordinate of your texture
	e.renderer.Copy(texture, nil, &messageRectangle)

	return nil
}
