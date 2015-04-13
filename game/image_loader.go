package game

type ImageLoader interface {
	LoadFile(path string) (Image, error)
}

type Image interface {
	Size() (w, h int)
	// src and dest must not be nil
	Draw(src, dest Rectangle)
}
