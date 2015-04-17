package game

import (
	"github.com/gonutz/strategy/images"
	"github.com/veandco/go-sdl2/sdl"
)

func NewGame(imageLoader images.ImageLoader, cam ScreenCamera) *Game {
	g := &Game{
		running: true,
	}
	g.state = NewPlayState(g, imageLoader, cam)
	return g
}

type Game struct {
	running   bool
	updating  bool
	state     GameState
	nextState GameState
}

func (g *Game) Running() bool {
	return g.running
}

func (g *Game) Quit() {
	g.running = false
}

func (g *Game) Update() {
	if g.state != nil {
		g.updating = true
		g.state.Update()
		g.updating = false
	}
	if g.nextState != nil {
		g.state = g.nextState
		g.nextState = nil
	}
}

func (g *Game) Draw() {
	if g.state != nil {
		g.state.Draw()
	}
}

func (g *Game) ChangeStateTo(state GameState) {
	g.nextState = state
}

func (g *Game) KeyDown(key sdl.Keycode) {
	if g.state != nil {
		g.state.KeyDown(key)
	}
}

func (g *Game) ScrollUp(x, y int) {
	if g.state != nil {
		g.state.ScrollUp(x, y)
	}
}

func (g *Game) ScrollDown(x, y int) {
	if g.state != nil {
		g.state.ScrollDown(x, y)
	}
}

func (g *Game) MouseMovedTo(x, y int) {
	if g.state != nil {
		g.state.MouseMovedTo(x, y)
	}
}
