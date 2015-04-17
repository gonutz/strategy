package images

type ImageSection struct {
	Image
	X, Y, W, H int
}

func (i *ImageSection) Size() (int, int) { return i.W, i.H }

func (i *ImageSection) Draw(src, dest Rect) {
	src.X += i.X
	src.Y += i.Y
	i.Image.Draw(src, dest)
}
