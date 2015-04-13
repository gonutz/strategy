package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

func newCamera(renderer *sdl.Renderer) *camera {
	return &camera{
		renderer: renderer,
	}
}

type camera struct {
	renderer           *sdl.Renderer
	renderTarget       *sdl.Texture
	screenW, screenH   int
	x, y               int
	displayW, displayH int
	worldW, worldH     int
	changed            bool
}

func (c *camera) setScreenSize(w, h int) {
	c.screenW, c.screenH = w, h
	c.changed = true
}

func (c *camera) SetVisibleHeight(h int) {
	if h < 1 {
		h = 1
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
		aspect := float32(c.screenW) / float32(c.screenH)
		c.displayW = int(float32(c.displayH)*aspect + 0.5)
		c.createRenderTarget()
		c.changed = false
	}
}

func (c *camera) createRenderTarget() {
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

func (c *camera) GetVisibleArea() (x, y, w, h int) {
	return c.x, c.y, c.displayW, c.displayH
}

func (c *camera) SetFocus(x, y int) {
	// TODO
}

func (c *camera) SetWorldSize(w, h int) {
	c.worldW, c.worldH = w, h
}

func (c *camera) Move(dx, dy int) {
	// TODO
}

// TODO eventually the zoom methods will be replaced with setting the visible
// height directly from inside the game state
func (c *camera) zoomIn() {
	c.SetVisibleHeight(c.displayH / 2)
}

func (c *camera) zoomOut() {
	c.SetVisibleHeight(c.displayH * 2)
}
