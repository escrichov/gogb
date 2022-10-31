package emulator

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

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
					var msg string
					if e.vsyncEnabled {
						msg = fmt.Sprintf("Vsync enabled")
					} else {
						msg = fmt.Sprintf("Vsync disabled")
					}
					e.SetMessage(msg, time.Second*3)
				case sdl.K_f:
					e.showFPS = !e.showFPS
				case sdl.K_p:
					snapshotFile := "snapshot.png"
					e.TakeSnapshot(snapshotFile)
					msg := fmt.Sprintf("Snapshot saved: %s", snapshotFile)
					e.SetMessage(msg, time.Second*3)
				case sdl.K_s:
					stateFile := "save.bess"
					e.BessStore(stateFile)
					msg := fmt.Sprintf("State saved: %s", stateFile)
					e.SetMessage(msg, time.Second*3)
				}
			}
		case *sdl.QuitEvent:
			e.stop = true
		}
	}
}
