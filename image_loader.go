package main

import (
	"github.com/gonutz/strategy/game"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

func newImageLoader(renderer *sdl.Renderer) *imageLoader {
	return &imageLoader{renderer, make(map[string]game.Image), nil}
}

type imageLoader struct {
	renderer *sdl.Renderer
	images   map[string]game.Image
	textures []*sdl.Texture
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
	l.textures = append(l.textures, texture)

	l.images[path] = &textureImage{l.renderer, texture}
	return l.images[path], nil
}

func (l *imageLoader) cleanUp() {
	for _, texture := range l.textures {
		texture.Destroy()
	}
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
