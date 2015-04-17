package images

type Rect struct {
	X, Y, W, H int
}

func (r Rect) TopLeft() (int, int) { return r.X, r.Y }
func (r Rect) Size() (int, int)    { return r.W, r.H }
