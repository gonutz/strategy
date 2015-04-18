package buildings

import "github.com/gonutz/strategy/images"

type Base struct {
	tiles  *images.TiledImage
	shadow images.Image
	x, y   int
	dead   bool
}

func NewBase(loader images.ImageLoader, x, y int) (*Base, error) {
	tileImage, err := loader.LoadImage("base.png")
	if err != nil {
		return nil, err
	}
	shadowImage, err := loader.LoadImage("base_shadow.png")
	if err != nil {
		return nil, err
	}
	return &Base{
		tiles: images.NewTiledImage(tileImage, images.TileImageDescription{
			TilesInX: 4,
			TilesInY: 2,
			Margin:   1,
			Spacing:  1,
		}),
		shadow: shadowImage,
		x:      x,
		y:      y,
	}, nil
}

func (b *Base) Draw(cam camera) {
	offsetX, offsetY, _, _ := cam.GetVisibleArea()
	const shadowX = 0
	const shadowY = 16
	w, h := b.shadow.Size()
	b.shadow.Draw(images.Rect{0, 0, w, h},
		images.Rect{b.x - offsetX + shadowX, b.y - offsetY + shadowY, w, h})

	tile := b.tiles.Tile(0)
	w, h = tile.Size()
	tile.Draw(images.Rect{0, 0, w, h},
		images.Rect{b.x - offsetX, b.y - offsetY, w, h})
}
