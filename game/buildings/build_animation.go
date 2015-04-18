package buildings

import (
	"github.com/gonutz/strategy/anim"
	"github.com/gonutz/strategy/images"
)

type BuildAnimation struct {
	animation *anim.TileAnimation
	x, y      int
}

func NewBuildAnimation(loader images.ImageLoader, x, y int) (*BuildAnimation, error) {
	tileImage, err := loader.LoadImage("build_13x5.png")
	if err != nil {
		return nil, err
	}
	desc := images.TileImageDescription{
		TilesInX: 13,
		TilesInY: 5,
		Margin:   1,
		Spacing:  1,
	}
	tiles := images.NewTiledImage(tileImage, desc)

	animation := anim.NewTileAnimation(
		anim.FrameDuration(3,
			anim.TileRange(tiles, 0, 35),
			anim.After(
				anim.Loop(5,
					anim.Sequence(
						anim.TileRange(tiles, 36, 38),
						anim.TileRange(tiles, 37, 37),
					),
				),
				func() { println("create tank") },
			),
			anim.TileRange(tiles, 39, 64),
		),
	)

	return &BuildAnimation{
		animation: animation,
		x:         x,
		y:         y,
	}, nil
}

func (b *BuildAnimation) Update() {
	b.animation.Update()
}

func (b *BuildAnimation) Draw(cam camera) {
	offsetX, offsetY, _, _ := cam.GetVisibleArea()
	b.animation.DrawAt(b.x-offsetX, b.y-offsetY)
}
