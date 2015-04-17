package main

import (
	"fmt"
	"github.com/gonutz/strategy/game"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

func main() {
	sdl.SetHint(sdl.HINT_RENDER_VSYNC, "1")

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	w, h := 640, 480
	flags := uint32(0)
	flags = sdl.WINDOW_FULLSCREEN // TODO debug code: use this to toggle fullscreen
	if flags == sdl.WINDOW_FULLSCREEN {
		w, h = screenResolution()
	} else {
		flags = sdl.WINDOW_RESIZABLE
	}
	window, renderer, err := sdl.CreateWindowAndRenderer(w, h, flags)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()
	defer renderer.Destroy()
	window.SetTitle("Strategy Game")
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	// TODO which one or should the user be able to choose one?
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "nearest")
	//sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")

	camera := newCamera(renderer)
	camera.setScreenSize(window.GetSize())
	camera.SetVisibleHeight(240)
	loader := newImageLoader(renderer)
	defer loader.cleanUp()
	g := game.NewGame(loader, camera)
	const frameIntervalInS = 1.0 / 60.0
	frameTimeInS := 0.0
	lastUpdate := time.Now()

	for g.Running() {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				g.Quit()
			case *sdl.KeyDownEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					g.Quit()
				} else {
					g.KeyDown(e.Keysym.Sym)
				}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					camera.setScreenSize(window.GetSize())
				}
			case *sdl.MouseMotionEvent:
				g.MouseMovedTo(int(e.X), int(e.Y))
			case *sdl.MouseWheelEvent:
				mouseX, mouseY, _ := sdl.GetMouseState()
				if e.Y > 0 {
					g.ScrollUp(mouseX, mouseY)
				}
				if e.Y < 0 {
					g.ScrollDown(mouseX, mouseY)
				}
			}
		}

		now := time.Now()
		frameTimeInS += now.Sub(lastUpdate).Seconds()
		lastUpdate = now
		if frameTimeInS > frameIntervalInS {
			for frameTimeInS > frameIntervalInS {
				g.Update()
				frameTimeInS -= frameIntervalInS
			}
			camera.startScene()
			g.Draw()
			camera.presentScene()
		}
	}
}

func screenResolution() (int, int) {
	var bounds sdl.Rect
	sdl.GetDisplayBounds(0, &bounds)
	return int(bounds.W), int(bounds.H)
}
