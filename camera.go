package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

func newCamera(renderer *sdl.Renderer) *camera {
	return &camera{renderer: renderer, changed: true}
}

type camera struct {
	renderer                 *sdl.Renderer
	renderTarget             *sdl.Texture
	screenW, screenH         int
	x, y                     int
	displayW, displayH       int
	maxDisplayW, maxDisplayH int
	changed                  bool
}

func (c *camera) setScreenSize(w, h int) {
	c.screenW, c.screenH = w, h
	c.changed = true
}

func (c *camera) SetMaximumScreenSize(w, h int) {
	c.maxDisplayW, c.maxDisplayH = w, h
}

func (c *camera) SetVisibleHeight(h int) {
	if h < 80 {
		h = 80
	}
	if h > 1600 {
		h = 1600
	}
	c.displayH = h
	c.changed = true
}

func (c *camera) startScene() {
	c.recalculate()
	c.renderer.SetRenderTarget(c.renderTarget)
	c.renderer.SetDrawColor(100, 200, 255, 255)
	c.renderer.Clear()
}

func (c *camera) recalculate() {
	if c.changed {
		if c.displayH > c.maxDisplayH {
			c.displayH = c.maxDisplayH
		}
		aspect := float32(c.screenW) / float32(c.screenH)
		c.displayW = int(float32(c.displayH)*aspect + 0.5)
		if c.displayW > c.maxDisplayW {
			c.displayW = c.maxDisplayW
			c.displayH = int(float32(c.displayW)/aspect + 0.5)
		}
		c.recreateRenderTarget()
		c.changed = false
	}
}

func (c *camera) recreateRenderTarget() {
	// TODO is SDL_PIXELFORMAT_UNKNOWN OK here?
	texture, err := c.renderer.CreateTexture(
		sdl.PIXELFORMAT_UNKNOWN,
		sdl.TEXTUREACCESS_TARGET,
		c.displayW,
		c.displayH,
	)
	if err != nil {
		fmt.Println("creating back buffer failed")
		return
	}
	c.renderTarget = texture
}

func (c *camera) presentScene() {
	c.renderer.SetRenderTarget(nil)
	c.renderer.Copy(c.renderTarget, nil, nil)
	c.renderer.Present()
}

func (c *camera) ScreenSize() (int, int) {
	return c.displayW, c.displayH
}

func (c *camera) ScreenToScene(x, y int) (int, int) {
	scale := float32(c.displayH) / float32(c.screenH)
	return int(float32(x)*scale + 0.5), int(float32(y)*scale + 0.5)
}

// TODO eventually the zoom methods may be replaced with setting the visible
// height directly from inside the game state
func (c *camera) ZoomIn() {
	c.SetVisibleHeight(int(float32(c.displayH)/1.1 + 0.5))
}

func (c *camera) ZoomOut() {
	c.SetVisibleHeight(int(float32(c.displayH)*1.1 + 0.5))
}
