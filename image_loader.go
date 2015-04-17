package main

import (
	"github.com/gonutz/strategy/images"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

func newImageLoader(renderer *sdl.Renderer) *imageLoader {
	return &imageLoader{renderer, make(map[string]images.Image), nil}
}

type imageLoader struct {
	renderer *sdl.Renderer
	images   map[string]images.Image
	textures []*sdl.Texture
}

func (l *imageLoader) LoadFile(path string) (images.Image, error) {
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

func (i *textureImage) Draw(src, dest images.Rect) {
	i.renderer.Copy(
		i.texture,
		&sdl.Rect{int32(src.X), int32(src.Y), int32(src.W), int32(src.H)},
		&sdl.Rect{int32(dest.X), int32(dest.Y), int32(dest.W), int32(dest.H)},
	)
}

func (i *textureImage) Size() (int, int) {
	_, _, w, h, _ := i.texture.Query()
	return w, h
}
