package images

type TiledImage struct {
	Image
	tiles []Image
}

type TileImageDescription struct {
	TilesInX int
	TilesInY int
	Margin   int
	Spacing  int
}

func NewTiledImage(image Image, desc TileImageDescription) *TiledImage {
	tileCount := desc.TilesInX * desc.TilesInY
	img := &TiledImage{image, make([]Image, tileCount)}
	w, h := image.Size()
	tileW := (w - 2*desc.Margin - (desc.TilesInX-1)*desc.Spacing) / desc.TilesInX
	tileH := (h - 2*desc.Margin - (desc.TilesInY-1)*desc.Spacing) / desc.TilesInY

	for i := 0; i < tileCount; i++ {
		tileX, tileY := i%desc.TilesInX, i/desc.TilesInX
		x := desc.Margin + tileX*(tileW+desc.Spacing)
		y := desc.Margin + tileY*(tileH+desc.Spacing)
		img.tiles[i] = &ImageSection{img.Image, x, y, tileW, tileH}
	}

	return img
}

func (img *TiledImage) Tile(i int) Image {
	return img.tiles[i]
}
