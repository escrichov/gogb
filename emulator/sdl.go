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

type Window struct {
	window   *sdl.Window
	surface  *sdl.Surface
	renderer *sdl.Renderer
	texture  *sdl.Texture
	font     *ttf.Font

	keyboardState []uint8
	frameBuffer   [WIDTH * HEIGHT]uint32

	// Vertical Sync (VSYNC) active
	vsyncEnabled bool

	// Window visible or only as a buffer
	showWindow bool

	// Window on fullscreen
	fullScreen bool

	// Frames per second
	frames                    uint64
	framesPerSecond           uint32
	framesCurrentSecond       uint32
	deltaTime                 uint64
	millisecondsPreviousFrame uint64
	showFPS                   bool

	// Console messages
	consoleMessage         string
	showMessage            bool
	consoleMessageDuration time.Duration
	consoleMessageStart    time.Time
}

func newWindow(title string, fontFilename string, windowScale float64, showWindow bool, vsyncAtStartup bool) (*Window, error) {
	window := &Window{
		vsyncEnabled: vsyncAtStartup,
		showFPS:      false,
		showMessage:  false,
		showWindow:   showWindow,
	}

	err := window.initializeSDL(title, fontFilename, windowScale, window.vsyncEnabled)
	if err != nil {
		return nil, err
	}

	blackColor := uint32(0)
	window.SetFramebufferColor(blackColor)

	return window, nil
}

func (w *Window) Destroy() {
	w.texture.Destroy()
	w.renderer.Destroy()
	w.window.Destroy()
	sdl.Quit()
}

func (w *Window) initializeSDL(windowName, fontFilename string, windowScale float64, vsyncEnable bool) error {
	var err error

	// SDL Initialization
	//var subsystemMask uint32 = sdl.INIT_VIDEO | sdl.INIT_AUDIO
	//if sdl.WasInit(subsystemMask) != subsystemMask {
	//	if err := sdl.Init(subsystemMask); err != nil {
	//		return err
	//	}
	//}
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	// Initialize SDL TTF
	if err := ttf.Init(); err != nil {
		return err
	}

	if w.showWindow {
		w.window, err = sdl.CreateWindow(
			windowName,
			sdl.WINDOWPOS_CENTERED,
			sdl.WINDOWPOS_CENTERED,
			int32(WIDTH*windowScale),
			int32(HEIGHT*windowScale),
			sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)
		if err != nil {
			return err
		}

		flags := uint32(sdl.RENDERER_ACCELERATED)
		if vsyncEnable {
			flags |= sdl.RENDERER_PRESENTVSYNC
		}
		w.renderer, err = sdl.CreateRenderer(w.window, -1, flags)
		if err != nil {
			return err
		}
	} else {
		w.surface, err = sdl.CreateRGBSurface(
			0,
			int32(WIDTH*windowScale),
			int32(HEIGHT*windowScale),
			32,
			0x00ff0000,
			0x0000ff00,
			0x000000ff,
			0xff000000)
		if err != nil {
			return err
		}
		w.renderer, err = sdl.CreateSoftwareRenderer(w.surface)
		if err != nil {
			return err
		}
	}

	// Creating a SDL texture that is used to display the color buffer
	w.texture, err = w.renderer.CreateTexture(
		uint32(sdl.PIXELFORMAT_RGBA32),
		sdl.TEXTUREACCESS_STREAMING,
		WIDTH,
		HEIGHT,
	)
	if err != nil {
		return err
	}

	// Point to Keyboard State
	w.keyboardState = sdl.GetKeyboardState()

	// Font
	w.font, err = ttf.OpenFont(fontFilename, 25)
	if err != nil {
		return err
	}

	// This keeps aspect ratio when resizing window
	w.renderer.SetLogicalSize(WIDTH, HEIGHT)

	// Vsync
	w.renderer.RenderSetVSync(w.vsyncEnabled)

	return nil
}

func (w *Window) calculateFPS() {
	w.frames++
	w.framesCurrentSecond++

	ticks := sdl.GetTicks64()
	w.deltaTime += ticks - w.millisecondsPreviousFrame
	w.millisecondsPreviousFrame = ticks
	if w.deltaTime >= 1000 {
		w.framesPerSecond = w.framesCurrentSecond
		w.framesCurrentSecond = 0
		w.deltaTime = 0
	}
}

func (w *Window) renderFrame() {
	w.calculateFPS()

	w.renderer.Clear()

	// Render FrameBuffer
	buf := unsafe.Pointer(&w.frameBuffer[0])
	framebufferBytes := unsafe.Slice((*byte)(buf), WIDTH*HEIGHT)
	w.texture.Update(nil, framebufferBytes, WIDTH*4)
	w.renderer.Copy(w.texture, nil, nil)

	if w.showFPS {
		w.drawFPS()
	}

	if w.showMessage {
		w.drawMessage()

		now := time.Now()
		if w.consoleMessageDuration != 0 && now.Sub(w.consoleMessageStart) > w.consoleMessageDuration {
			w.showMessage = false
		}
	}

	w.renderer.Present()
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

func (w *Window) TakeSnapshot(filename string) error {
	width, height, err := w.renderer.GetOutputSize()
	if err != nil {
		return err
	}

	surface, err := sdl.CreateRGBSurface(0, width, height, 32, 0x00ff0000, 0x0000ff00, 0x000000ff, 0xff000000)
	if err != nil {
		return err
	}

	err = w.renderer.ReadPixels(nil, sdl.PIXELFORMAT_ARGB8888, surface.Data(), int(surface.Pitch))
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

func (w *Window) drawFPS() error {
	color := sdl.Color{R: 0, G: 255, B: 0, A: 255}
	framesString := fmt.Sprintf("FPS %d", w.framesPerSecond)
	return w.drawText(WIDTH*4-100, 0, framesString, color)
}

func (w *Window) drawMessage() error {
	color := sdl.Color{R: 0, G: 0, B: 255, A: 255}
	return w.drawText(10, HEIGHT*4-35, w.consoleMessage, color)
}

func (w *Window) SetMessage(message string, duration time.Duration) {
	w.showMessage = true
	w.consoleMessage = message
	w.consoleMessageDuration = duration
	w.consoleMessageStart = time.Now()
}

func (w *Window) drawText(x, y int32, text string, color sdl.Color) error {

	// as TTF_RenderText_Solid could only be used on
	// SDL_Surface then you have to create the surface first
	surface, err := w.font.RenderUTF8Solid(text, color)
	if err != nil {
		return err
	}
	defer surface.Free()

	// now you can convert it into a texture
	texture, err := w.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()

	textWidth, textHeight, err := w.font.SizeUTF8(text)

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
	w.renderer.Copy(texture, nil, &messageRectangle)

	return nil
}

func (w *Window) SetFramebufferColor(color uint32) {
	for i, _ := range w.frameBuffer {
		w.frameBuffer[i] = color
	}
}

func (w *Window) ToggleVsync() {
	w.SetVsync(!w.vsyncEnabled)
}

func (w *Window) SetVsync(active bool) {
	w.vsyncEnabled = active
	w.renderer.RenderSetVSync(w.vsyncEnabled)
	if w.vsyncEnabled {
		sdl.GLSetSwapInterval(1)
	} else {
		sdl.GLSetSwapInterval(0)
	}
}

func (w *Window) GetVsync() bool {
	return w.vsyncEnabled
}

func (w *Window) ToggleShowFPS() {
	w.showFPS = !w.showFPS
}

func (w *Window) GetFramebuffer() *[WIDTH * HEIGHT]uint32 {
	return &w.frameBuffer
}

func (w *Window) GetKeyboardState() []uint8 {
	return w.keyboardState
}

func (w *Window) ToggleFullScreen() {
	w.fullScreen = !w.fullScreen
	if w.fullScreen {
		w.window.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	} else {
		w.window.SetFullscreen(0)
		w.window.SetPosition(sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED)
	}
}
