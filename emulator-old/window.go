package emulator_old

import (
	"emulator-go/emulator/gb/utils"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"unsafe"
)
import "github.com/veandco/go-sdl2/ttf"

func SDLInit() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	// Initialize SDL TTF
	if err := ttf.Init(); err != nil {
		return err
	}

	return nil
}

type Window struct {
	fps                       uint
	millisecondsPreviousFrame uint64
	millisecondsPerFrame      uint64
	deltaTime                 uint64
	window                    *sdl.Window
	renderer                  *sdl.Renderer
	texture                   *sdl.Texture
	font                      *ttf.Font
	colorBuffer               []uint32
	WindowWidth               int32
	WindowHeight              int32
	showFPS                   bool
	showGrid                  bool
}

func (e *Window) Initialize(title string, width int32, height int32, fps uint, showFPS bool, showGrid bool) error {
	var err error

	e.WindowWidth = width
	e.WindowHeight = height
	e.fps = fps
	e.showFPS = showFPS
	e.showGrid = showGrid

	e.millisecondsPreviousFrame = 0
	e.deltaTime = 1
	e.millisecondsPerFrame = uint64(1000 / e.fps)
	e.colorBuffer = make([]uint32, e.WindowWidth*e.WindowHeight)

	e.window, err = sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		e.WindowWidth*2,
		e.WindowHeight*2,
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
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		e.WindowWidth,
		e.WindowHeight,
	)
	if err != nil {
		return err
	}

	e.font, err = ttf.OpenFont("assets/fonts/arial.ttf", 25)
	if err != nil {
		return err
	}

	return nil
}

func (e *Window) Destroy() error {
	if e.texture != nil {
		e.texture.Destroy()
	}
	if e.renderer != nil {
		e.renderer.Destroy()
	}
	if e.window != nil {
		e.window.Destroy()
	}
	if e.font != nil {
		e.font.Close()
	}
	sdl.Quit()
	return nil
}

func (e *Window) SetPixel(x int32, y int32, color uint32) error {
	if x >= 0 && y >= 0 && x < e.WindowWidth && y < e.WindowHeight {
		e.colorBuffer[utils.RowColtoPos(int(y), int(x), int(e.WindowWidth))] = color
		return nil
	}

	return fmt.Errorf("invalid position: %d, %d", x, y)
}

func (e *Window) SetColorBuffer(colorBuffer []uint32) error {
	if len(colorBuffer) == len(e.colorBuffer) {
		copy(e.colorBuffer, colorBuffer)
		return nil
	}

	return fmt.Errorf("invalid buffer length: %d", len(colorBuffer))
}

func (e *Window) clearColorBuffer(color uint32) {
	for y := int32(0); y < e.WindowHeight; y++ {
		for x := int32(0); x < e.WindowWidth; x++ {
			e.colorBuffer[utils.RowColtoPos(int(y), int(x), int(e.WindowWidth))] = color
		}
	}
}

func (e *Window) drawGrid() {
	for j := int32(0); j < e.WindowHeight; j++ {
		for i := int32(0); i < e.WindowWidth; i++ {
			if i%8 == 0 || j%8 == 0 {
				e.colorBuffer[utils.RowColtoPos(int(j), int(i), int(e.WindowWidth))] = 0xFF00FF00
			}
		}
	}
}

func (e *Window) drawFPS() error {
	// this is the color in rgb format,
	// maxing out all would give you the color white,
	// and it will be your text's color
	whiteColor := sdl.Color{R: 255, G: 255, B: 255, A: 255}

	// as TTF_RenderText_Solid could only be used on
	// SDL_Surface then you have to create the surface first
	var fpsString string

	if e.deltaTime != 0 {
		fpsString = fmt.Sprintf("%.2f FPS", float64(1000.0/e.deltaTime))
	} else {
		fpsString = fmt.Sprintf("%d FPS", e.fps)
	}

	surface, err := e.font.RenderUTF8Solid(fpsString, whiteColor)
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

	messageRectangle := sdl.Rect{
		X: e.WindowWidth*2 - 101, // Controls the rect's x coordinate
		Y: 0,                     // controls the rect's y coordinate
		W: 100,                   // controls the width of the rect
		H: 50,                    // controls the height of the rect
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

func (e *Window) Update() {
	timeToWait := e.millisecondsPerFrame - (sdl.GetTicks64() - e.millisecondsPreviousFrame)
	if timeToWait > 0 && timeToWait <= e.millisecondsPerFrame {
		sdl.Delay(uint32(timeToWait))
	}

	ticks := sdl.GetTicks64()
	e.deltaTime = ticks - e.millisecondsPreviousFrame
	e.millisecondsPreviousFrame = ticks
}

func (e *Window) drawColorBuffer() {
	buf := unsafe.Pointer(&e.colorBuffer[0])
	pixels := unsafe.Slice((*byte)(buf), e.WindowWidth*e.WindowHeight)
	e.texture.Update(nil, pixels, int(e.WindowWidth*4))
	e.renderer.Copy(e.texture, nil, nil)
}

func (e *Window) Render() {
	e.renderer.Clear()

	// Draw
	if e.showGrid {
		e.drawGrid()
	}
	e.drawColorBuffer()
	if e.showFPS {
		e.drawFPS()
	}

	e.renderer.Present()
}
