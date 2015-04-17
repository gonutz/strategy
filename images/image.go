package images

// Image represents a rectangular pixel image that can be drawn. The top-left is
// defined to be at 0,0 and the size in pixels can be queried with Size.
type Image interface {

	// Size reports the Image's width and height in pixels.
	Size() (width, height int)

	// Draw draws the src portion of the image at the dest area. The space of
	// dest is determined by the concrete Image implementation. It could be
	// screen space but it could also be world space or some other transformed
	// space.
	// The parameters src and dest must not be nil.
	// The source parameter should not lie outside the Image's boundaries, the
	// behavior is otherwise undefined.
	Draw(src, dest Rect)
}
