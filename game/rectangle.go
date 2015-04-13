package game

type Rectangle interface {
	TopLeft() (x, y int)
	Size() (w, h int)
}

type rect struct {
	x, y, w, h int
}

func (r rect) TopLeft() (int, int) { return r.x, r.y }
func (r rect) Size() (int, int)    { return r.w, r.h }
