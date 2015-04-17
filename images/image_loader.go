package images

type ImageLoader interface {
	LoadFile(path string) (Image, error)
}
