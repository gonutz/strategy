package game

import "github.com/veandco/go-sdl2/sdl"

type GameState interface {
	Update()
	Draw()
	KeyDown(key sdl.Keycode)
	ScrollUp(x, y int)
	ScrollDown(x, y int)
	MouseMovedTo(x, y int)
}
