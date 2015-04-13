package game

type SceneCamera interface {
	// GetVisibleArea returns the world-space rectangle that is currently visible.
	GetVisibleArea() (x, y, w, h int)

	// SetWorldSize sets the size of the world. The camera is not allowed to
	// move outside the boundaries (0, 0, w, h). If the world size is smaller
	// than the visible rectangle it will be fully visible.
	// TODO or clamp the visible height to world coordinates?
	SetWorldSize(w, h int)

	// SetVisibleHeight sets the height in world-space of the camera's visible
	// rectangle. The visible width is automatically.
	SetVisibleHeight(h int)

	// SetFocus sets the desired center point of the camera in world-space
	// coordinates. This point might not end up in the screen center if the
	//camera is clamped to the world size borders but it will be always visible.
	SetFocus(x, y int)

	// Move moves the visible camera rectangle by the given amount. If the
	// camera leaves the world boundaries it is clamped.
	Move(dx, dy int)
}
