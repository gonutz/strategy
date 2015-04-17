package images

type ImageSection struct {
	Image
	X, Y, W, H int
}

func (i *ImageSection) Size() (int, int) { return i.W, i.H }

// TODO use Rect instead of Rectangle and make all calls use the concrete type.
// There is probably no need for the indirection and using Rect directly (not
// *Rect) this could just manipulate X and Y and forward the call then. No
// additional offsetRectangle would be needed that way.
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
