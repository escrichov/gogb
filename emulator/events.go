package emulator

import "github.com/veandco/go-sdl2/sdl"

func (e *Emulator) manageKeyboardEvents() {
	// Manage Keyboard Events
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.KeyboardEvent:
			keyboardEvent, ok := event.(*sdl.KeyboardEvent)
			if !ok {
				continue
			}
			if keyboardEvent.Type == sdl.KEYDOWN {
				switch keyboardEvent.Keysym.Sym {
				case sdl.K_ESCAPE:
					e.stop = true
				case sdl.K_v:
					e.vsyncEnabled = !e.vsyncEnabled
					e.renderer.RenderSetVSync(e.vsyncEnabled)
				case sdl.K_s:
					e.BessStore("save.bess")
				}
			}
		case *sdl.QuitEvent:
			e.stop = true
		}
	}
}
