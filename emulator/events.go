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
					e.window.ToggleVsync()
					var msg string
					if e.window.GetVsync() {
						msg = fmt.Sprintf("Vsync enabled")
					} else {
						msg = fmt.Sprintf("Vsync disabled")
					}
					e.window.SetMessage(msg, time.Second*3)
				case sdl.K_t:
					e.window.ToggleShowFPS()
				case sdl.K_f:
					e.window.ToggleFullScreen()
				case sdl.K_k:
					snapshotFile := "snapshot.png"
					e.window.TakeSnapshot(snapshotFile)
					msg := fmt.Sprintf("Snapshot saved: %s", snapshotFile)
					e.window.SetMessage(msg, time.Second*3)
				case sdl.K_s:
					stateFile := "save.bess"
					e.BessStore(stateFile)
					msg := fmt.Sprintf("State saved: %s", stateFile)
					e.window.SetMessage(msg, time.Second*3)
				case sdl.K_r:
					e.reset = true
					e.window.SetMessage("Reset", time.Second*3)
				case sdl.K_p:
					e.pause = !e.pause
					if e.pause {
						e.window.SetMessage("Paused", 0)
					} else {
						e.window.SetMessage("Continue", time.Second*3)
					}

				}
			}
		case *sdl.QuitEvent:
			e.stop = true
		}
	}
}
