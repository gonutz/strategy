package images

type TiledImage struct {
	Image
	tiles []Image
}

func NewTiledImageWithTileCounts(image Image, tilesInX, tilesInY int) *TiledImage {
	w, h := image.Size()
	tileW := w / tilesInX
	tileH := h / tilesInY
	img := &TiledImage{image, nil}
	img.computeTiles(tileW, tileH)
	return img
}

func (img *TiledImage) computeTiles(tileW, tileH int) {
	w, h := img.Image.Size()
	tilesInX, tilesInY := w/tileW, h/tileH
	tileCount := tilesInX * tilesInY
	img.tiles = make([]Image, tileCount)
	for i := 0; i < tileCount; i++ {
		tileX, tileY := i%tilesInX, i/tilesInX
		x, y := tileX*tileW, tileY*tileH
		img.tiles[i] = &ImageSection{img.Image, x + 1, y + 1, tileW - 2, tileH - 2}
	}
}

func (img *TiledImage) Tile(i int) Image {
	return img.tiles[i]
}
