package emulator_old

import "github.com/veandco/go-sdl2/sdl"

func (e *Emulator) ProcessEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			e.isRunning = false
		case *sdl.WindowEvent:
			windowEvent, ok := event.(*sdl.WindowEvent)
			if !ok {
				continue
			}
			if windowEvent.Event == sdl.WINDOWEVENT_CLOSE {
				e.isRunning = false
			}
		case *sdl.KeyboardEvent:
			keyboardEvent, ok := event.(*sdl.KeyboardEvent)
			if !ok {
				continue
			}
			if keyboardEvent.Type == sdl.KEYDOWN {
				switch keyboardEvent.Keysym.Sym {
				case sdl.K_ESCAPE:
					e.isRunning = false
				case sdl.K_f:
					e.MainWindow.showFPS = !e.MainWindow.showFPS
					e.BackgroundWindow.showFPS = !e.BackgroundWindow.showFPS
				case sdl.K_g:
					e.MainWindow.showGrid = !e.MainWindow.showGrid
					e.BackgroundWindow.showGrid = !e.BackgroundWindow.showGrid
				}
			}
		}
	}
}
