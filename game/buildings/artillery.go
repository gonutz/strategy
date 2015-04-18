package buildings

import "github.com/gonutz/strategy/images"

type Artillery struct {
	tiles       *images.TiledImage
	shadow      *images.TiledImage
	rotationDeg int
	x, y        int
	shooting    int
	dead        bool
}

func NewArtillery(loader images.ImageLoader, x, y int) (*Artillery, error) {
	tileImage, err := loader.LoadImage("artillery.png")
	if err != nil {
		return nil, err
	}
	return &Artillery{
		tiles: images.NewTiledImage(tileImage, images.TileImageDescription{
			TilesInX: 6,
			TilesInY: 3,
			Margin:   1,
			Spacing:  1,
		}),
		x: x,
		y: y,
	}, nil
}

func (a *Artillery) RotateBy(degrees int) {
	a.rotationDeg += degrees
}

func (a *Artillery) Update() {
	a.rotationDeg++
}

func (a *Artillery) Draw(cam camera) {
	var tileIndex int
	if a.dead {
		tileIndex = 16
	} else {
		rot := (normalizeRotation(a.rotationDeg)*2 + 45) % 720
		tileIndex = (rot / 90) * 2
		if a.shooting > 0 {
			tileIndex++
		}
	}
	tile := a.tiles.Tile(tileIndex)
	w, h := tile.Size()
	offsetX, offsetY, _, _ := cam.GetVisibleArea()
	tile.Draw(
		images.Rect{0, 0, w, h},
		images.Rect{a.x - offsetX, a.y - offsetY, w, h},
	)
}

type camera interface {
	GetVisibleArea() (x, y, w, h int)
}

func normalizeRotation(degrees int) int {
	return (degrees%360 + 360) % 360
}
