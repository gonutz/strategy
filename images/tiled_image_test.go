package images

import "testing"

func TestTilesAreEquallySizedImageSections_IndexedTopToBottom_LeftToRight(t *testing.T) {
	spy := newSpyImage(20, 60)
	tiled := NewTiledImageWithTileCounts(spy, TileImageDescription{
		TilesInX: 2,
		TilesInY: 3,
		Margin:   0,
		Spacing:  0,
	})

	// the tile indices are as follows:
	// 01
	// 23
	// 45
	checkTile(t, tiled.Tile(0), spy, 0, 0, 10, 20)
	checkTile(t, tiled.Tile(1), spy, 10, 0, 10, 20)
	checkTile(t, tiled.Tile(2), spy, 0, 20, 10, 20)
	checkTile(t, tiled.Tile(3), spy, 10, 20, 10, 20)
	checkTile(t, tiled.Tile(4), spy, 0, 40, 10, 20)
	checkTile(t, tiled.Tile(5), spy, 10, 40, 10, 20)
}

func TestTileMarginIsBorderAroundWholeTileImage(t *testing.T) {
	spy := newSpyImage(24, 64)
	tiled := NewTiledImageWithTileCounts(spy, TileImageDescription{
		TilesInX: 2,
		TilesInY: 3,
		Margin:   2,
		Spacing:  0,
	})

	checkTile(t, tiled.Tile(0), spy, 2, 2, 10, 20)
	checkTile(t, tiled.Tile(1), spy, 12, 2, 10, 20)
	checkTile(t, tiled.Tile(2), spy, 2, 22, 10, 20)
	checkTile(t, tiled.Tile(3), spy, 12, 22, 10, 20)
	checkTile(t, tiled.Tile(4), spy, 2, 42, 10, 20)
	checkTile(t, tiled.Tile(5), spy, 12, 42, 10, 20)
}

func TestTileSpacingIsThePixelCountBetweenAdjacentTiles(t *testing.T) {
	spy := newSpyImage(23, 66)
	tiled := NewTiledImageWithTileCounts(spy, TileImageDescription{
		TilesInX: 2,
		TilesInY: 3,
		Margin:   0,
		Spacing:  3,
	})

	checkTile(t, tiled.Tile(0), spy, 0, 0, 10, 20)
	checkTile(t, tiled.Tile(1), spy, 13, 0, 10, 20)
	checkTile(t, tiled.Tile(2), spy, 0, 23, 10, 20)
	checkTile(t, tiled.Tile(3), spy, 13, 23, 10, 20)
	checkTile(t, tiled.Tile(4), spy, 0, 46, 10, 20)
	checkTile(t, tiled.Tile(5), spy, 13, 46, 10, 20)
}

func TestTilesCanHaveMarginAndSpacing(t *testing.T) {
	spy := newSpyImage(25, 68)
	tiled := NewTiledImageWithTileCounts(spy, TileImageDescription{
		TilesInX: 2,
		TilesInY: 3,
		Margin:   1,
		Spacing:  3,
	})

	checkTile(t, tiled.Tile(0), spy, 1, 1, 10, 20)
	checkTile(t, tiled.Tile(1), spy, 14, 1, 10, 20)
	checkTile(t, tiled.Tile(2), spy, 1, 24, 10, 20)
	checkTile(t, tiled.Tile(3), spy, 14, 24, 10, 20)
	checkTile(t, tiled.Tile(4), spy, 1, 47, 10, 20)
	checkTile(t, tiled.Tile(5), spy, 14, 47, 10, 20)
}

func checkTile(t *testing.T, tileImage Image, spy *spyImage, x, y, w, h int) {
	if gotW, gotH := tileImage.Size(); gotW != w || gotH != h {
		t.Errorf("size not %vx%v but %vx%v", w, h, gotW, gotH)
	}
	tileImage.Draw(Rect{}, Rect{})
	if gotX, gotY := spy.drawSrc.TopLeft(); gotX != x || gotY != y {
		t.Errorf("top-left not at %v,%v but %v,%v", x, y, gotX, gotY)
	}
}

func newSpyImage(w, h int) *spyImage {
	return &spyImage{w: w, h: h}
}

type spyImage struct {
	w, h     int
	drawSrc  Rect
	drawDest Rect
}

func (i *spyImage) Size() (w, h int) {
	return i.w, i.h
}

func (i *spyImage) Draw(src, dest Rect) {
	i.drawSrc = src
	i.drawDest = dest
}
