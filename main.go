package main

import (
	"fmt"
	"github.com/gonutz/strategy/game"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
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
	flags |= sdl.WINDOW_FULLSCREEN // TODO debug code: use this to toggle fullscreen
	if flags&sdl.WINDOW_FULLSCREEN != 0 {
		w, h = screenResolution()
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
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "nearest")

	camera := newCamera(renderer)
	camera.setScreenSize(window.GetSize())
	camera.SetVisibleHeight(480)
	g := game.NewGame(newImageLoader(renderer), camera)
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
				}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					camera.setScreenSize(window.GetSize())
				}
			case *sdl.MouseWheelEvent:
				if e.Y > 0 {
					camera.zoomIn()
				}
				if e.Y < 0 {
					camera.zoomOut()
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

func newImageLoader(renderer *sdl.Renderer) game.ImageLoader {
	return &imageLoader{renderer, make(map[string]game.Image)}
}

type imageLoader struct {
	renderer *sdl.Renderer
	images   map[string]game.Image
}

func (l *imageLoader) LoadFile(path string) (game.Image, error) {
	if img, ok := l.images[path]; ok {
		return img, nil
	}

	surface, err := img.Load(path)
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	texture, err := l.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}

	l.images[path] = &textureImage{l.renderer, texture}
	return l.images[path], nil
}

type textureImage struct {
	renderer *sdl.Renderer
	texture  *sdl.Texture
}

func (i *textureImage) Draw(src, dest game.Rectangle) {
	sx, sy := src.TopLeft()
	sw, sh := src.Size()
	dx, dy := dest.TopLeft()
	dw, dh := dest.Size()
	i.renderer.Copy(
		i.texture,
		&sdl.Rect{int32(sx), int32(sy), int32(sw), int32(sh)},
		&sdl.Rect{int32(dx), int32(dy), int32(dw), int32(dh)},
	)
}

func (i *textureImage) Size() (int, int) {
	_, _, w, h, _ := i.texture.Query()
	return w, h
}
