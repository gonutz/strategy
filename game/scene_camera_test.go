package game

import "testing"

type screenSize struct{ w, h int }

func (s screenSize) ScreenSize() (int, int)            { return s.w, s.h }
func (s screenSize) ZoomIn()                           {}
func (s screenSize) ZoomOut()                          {}
func (s screenSize) ScreenToScene(x, y int) (int, int) { return x, y }
func (s screenSize) SetMaximumScreenSize(w, h int)     {}

func TestFocussingTopLeftMovesVisibleRectThere(t *testing.T) {
	cam := newSceneCamera(screenSize{800, 600})
	cam.SetWorldSize(2000, 2000)
	cam.SetFocus(0, 0)

	x, y, w, h := cam.GetVisibleArea()

	if x != 0 || y != 0 || w != 800 || h != 600 {
		t.Errorf("visible area was %v %v %v %v", x, y, w, h)
	}
}

func TestFocussingBottomRightMovesVisibleRectThere(t *testing.T) {
	cam := newSceneCamera(screenSize{800, 600})
	cam.SetWorldSize(2000, 1000)
	cam.SetFocus(2000, 1000)

	x, y, w, h := cam.GetVisibleArea()

	if x != 2000-800 || y != 1000-600 || w != 800 || h != 600 {
		t.Errorf("visible area was %v %v %v %v", x, y, w, h)
	}
}

func TestFocussingTopRightMovesVisibleRectThere(t *testing.T) {
	cam := newSceneCamera(screenSize{800, 600})
	cam.SetWorldSize(2000, 1000)
	cam.SetFocus(2000, 0)

	x, y, w, h := cam.GetVisibleArea()

	if x != 2000-800 || y != 0 || w != 800 || h != 600 {
		t.Errorf("visible area was %v %v %v %v", x, y, w, h)
	}
}

func TestFocussingBottomLeftMovesVisibleRectThere(t *testing.T) {
	cam := newSceneCamera(screenSize{800, 600})
	cam.SetWorldSize(2000, 1000)
	cam.SetFocus(0, 1000)

	x, y, w, h := cam.GetVisibleArea()

	if x != 0 || y != 1000-600 || w != 800 || h != 600 {
		t.Errorf("visible area was %v %v %v %v", x, y, w, h)
	}
}

func TestFocusPointIsInCenterOfVisibleRect(t *testing.T) {
	cam := newSceneCamera(screenSize{800, 599})
	cam.SetWorldSize(2000, 2000)
	cam.SetFocus(700, 500)

	x, y, w, h := cam.GetVisibleArea()

	if x != 700-400 || y != 500-599/2 || w != 800 || h != 599 {
		t.Errorf("visible area was %v %v %v %v", x, y, w, h)
	}
}

func TestMovingCameraMovesVisibleRectangle(t *testing.T) {
	cam := newSceneCamera(screenSize{800, 600})
	cam.SetWorldSize(2000, 2000)
	cam.SetFocus(100, 100)
	cam.Move(-50, -25) // visible rectangle should stay in top-left corner
	cam.Move(100, 200) // camera should move right 100 and down 200

	x, y, w, h := cam.GetVisibleArea()

	if x != 100 || y != 200 || w != 800 || h != 600 {
		t.Errorf("visible area was %v %v %v %v", x, y, w, h)
	}
}
