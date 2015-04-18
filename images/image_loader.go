package images

type ImageLoader interface {
	LoadImage(ID string) (Image, error)
}
