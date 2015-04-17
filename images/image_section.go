package images

type ImageSection struct {
	Image
	X, Y, W, H int
}

func (i *ImageSection) Size() (int, int) { return i.W, i.H }

func (i *ImageSection) Draw(src, dest Rectangle) {
	i.Image.Draw(&offsetRectangle{src, i.X, i.Y}, dest)
}

type offsetRectangle struct {
	Rectangle
	offsetX, offsetY int
}

func (r *offsetRectangle) TopLeft() (int, int) {
	x, y := r.Rectangle.TopLeft()
	return x + r.offsetX, y + r.offsetY
}
