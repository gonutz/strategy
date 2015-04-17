package images

type Image interface {
	Size() (w, h int)
	// src and dest must not be nil
	Draw(src, dest Rectangle)
}
