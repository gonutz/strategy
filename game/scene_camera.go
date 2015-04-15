package game

type SceneCamera interface {
	// SetWorldSize sets the size of the world. The camera is not allowed to
	// move outside the boundaries (0, 0, w, h). If the world size is smaller
	// than the visible rectangle it will be fully visible.
	// TODO or clamp the visible height to world coordinates?
	SetWorldSize(w, h int)

	// SetVisibleHeight sets the height in world-space of the camera's visible
	// rectangle. The visible width is calculated automatically from the
	// screen's aspect ratio.
	//SetVisibleHeight(h int)

	// SetFocus sets the desired center point of the camera in world-space
	// coordinates. This point might not end up in the screen center if the
	//camera is clamped to the world size borders but it will always be visible.
	SetFocus(x, y int)

	// Move moves the visible camera rectangle by the given amount. If the
	// camera leaves the world boundaries it is clamped.
	Move(dx, dy int)

	// GetVisibleArea returns the world-space rectangle that is currently visible.
	GetVisibleArea() (x, y, w, h int)

	// ScreenToWorld transforms the given screen-space coordinates into
	// world-space coordinates
	ScreenToWorld(x, y int) (int, int)

	// ZoomIn maginifies the currently visible rectangle.
	ZoomIn(mouseX, mouseY int)

	// ZoomOut minifies the currently visible rectangle.
	ZoomOut(mouseX, mouseY int)
}

func newSceneCamera(cam ScreenCamera) SceneCamera {
	return &sceneCamera{screenCamera: cam}
}

type ScreenCamera interface {
	// SetMaximumScreenSize defines the maximum width and height that the camera
	// might show at any time. Zooming out will not be possible beyond these
	// boundaries.
	SetMaximumScreenSize(w, h int)
	ScreenSize() (w, h int)
	ScreenToScene(x, y int) (int, int)
	ZoomIn()
	ZoomOut()
}

type sceneCamera struct {
	screenCamera   ScreenCamera
	worldW, worldH int
	focusX, focusY int
}

func (c *sceneCamera) SetWorldSize(w, h int) {
	c.worldW, c.worldH = w, h
	c.screenCamera.SetMaximumScreenSize(w, h)
}

func (c *sceneCamera) SetFocus(x, y int) {
	c.focusX, c.focusY = x, y
}

func (c *sceneCamera) Move(dx, dy int) {
	x, y, w, h := c.GetVisibleArea()
	x += dx
	y += dy
	c.focusX, c.focusY = x+w/2, y+h/2
}

func (c *sceneCamera) GetVisibleArea() (x, y, w, h int) {
	w, h = c.screenCamera.ScreenSize()
	x, y = c.focusX-w/2, c.focusY-h/2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	maxX := c.worldW - w
	if x > maxX {
		x = maxX
	}
	maxY := c.worldH - h
	if y > maxY {
		y = maxY
	}
	return
}

func (c *sceneCamera) ScreenToWorld(x, y int) (int, int) {
	x, y = c.screenCamera.ScreenToScene(x, y)
	offsetX, offsetY, _, _ := c.GetVisibleArea()
	return x + offsetX, y + offsetY
}

func (c *sceneCamera) ZoomIn(mouseX, mouseY int) {
	c.zoomAt(mouseX, mouseY, c.screenCamera.ZoomIn)
}

func (c *sceneCamera) ZoomOut(mouseX, mouseY int) {
	c.zoomAt(mouseX, mouseY, c.screenCamera.ZoomOut)
}

func (c *sceneCamera) zoomAt(mouseX, mouseY int, zoomFunc func()) {
	x, y := c.ScreenToWorld(mouseX, mouseY)
	zoomFunc()
	zoomedX, zoomedY := c.ScreenToWorld(mouseX, mouseY)
	offsetX, offsetY := x-zoomedX, y-zoomedY
	c.Move(offsetX, offsetY)
}
