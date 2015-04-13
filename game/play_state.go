package game

import (
	"github.com/mewmew/tmx"
	"os"
	"path/filepath"
)

func NewPlayState(context GameStateChanger, imageLoader ImageLoader, camera SceneCamera) *PlayState {
	state := &PlayState{context: context, camera: camera}
	state.loadMap(imageLoader, "crater_world.tmx")
	return state
}

type PlayState struct {
	context GameStateChanger
	camera  SceneCamera
	tileMap *tileMap
}

type GameStateChanger interface {
	ChangeStateTo(state GameState)
}

func (s *PlayState) Update() {}

func (s *PlayState) Draw() {
	w, h := s.tileMap.size()
	s.tileMap.draw(0, 0, w, h)
}

func (s *PlayState) loadMap(imageLoader ImageLoader, ID string) {
	resourcePath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "gonutz", "strategy", "rsc")
	mapPath := filepath.Join(resourcePath, ID)
	m, err := tmx.Open(mapPath)
	if err != nil {
		panic(err)
	}

	s.tileMap = &tileMap{
		tileW: m.TileWidth,
		tileH: m.TileHeight,
		width: m.Width,
		tiles: make([]tile, m.Width*m.Height),
	}

	var tileImages []*tileImage
	for _, tileSet := range m.Tilesets {
		img, err := newTileImage(
			imageLoader,
			filepath.Join(resourcePath, tileSet.Image.Source),
			tileSet.TileWidth, tileSet.TileHeight,
			tileSet.Margin, tileSet.Spacing,
			tileSet.FirstGID)
		if err != nil {
			panic(err)
		}
		tileImages = append(tileImages, img)
	}

	ground := m.Layers[0] // TODO look up the Layer with name "ground"
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			s.tileMap.tileAt(x, y).image = lookUpTileImage(tileImages, ground.GetGID(x, y))
		}
	}
}

func lookUpTileImage(images []*tileImage, gid int) Image {
	for _, tileImage := range images {
		if img := tileImage.imageWithID(gid); img != nil {
			return img
		}
	}
	return nil
}

type tileMap struct {
	tileW, tileH int
	width        int
	tiles        []tile
}

func (m *tileMap) size() (int, int) {
	return m.width, len(m.tiles) / m.width
}

func (m *tileMap) draw(left, top, width, height int) {
	for x := left; x < left+width; x++ {
		for y := top; y < top+height; y++ {
			m.tileAt(x, y).image.Draw(
				rect{0, 0, m.tileW, m.tileH},
				rect{x * m.tileW, y * m.tileH, m.tileW, m.tileH},
			)
		}
	}
}

func (m *tileMap) tileAt(x, y int) *tile {
	return &m.tiles[x+y*m.width]
}

type tile struct {
	image Image
}

func newTileImage(imageLoader ImageLoader, path string, tileW, tileH, margin, spacing, firstID int) (*tileImage, error) {
	img, err := imageLoader.LoadFile(path)
	if err != nil {
		return nil, err
	}
	w, h := img.Size()
	tilesX := (w - margin) / (tileW + spacing)
	tilesY := (h - margin) / (tileH + spacing)
	lastID := firstID + tilesX*tilesY - 1
	tileImage := &tileImage{img, tileW, tileH, margin, spacing, firstID, lastID, tilesX, tilesY}
	return tileImage, nil
}

type tileImage struct {
	image        Image
	tileW, tileH int
	margin       int
	spacing      int
	firstID      int

	lastID                 int
	tileCountX, tileCountY int
}

func (i *tileImage) imageWithID(id int) Image {
	if id < i.firstID || id > i.lastID {
		return nil
	}
	index := id - i.firstID
	x, y := index%i.tileCountX, index/i.tileCountX
	imgX := i.margin + x*(i.tileW+i.spacing)
	imgY := i.margin + y*(i.tileH+i.spacing)
	return &imageSection{i.image, imgX, imgY, i.tileW, i.tileH}
}

type imageSection struct {
	Image
	x, y, w, h int
}

func (i *imageSection) Size() (int, int) { return i.w, i.h }

func (i *imageSection) Draw(src, dest Rectangle) {
	i.Image.Draw(&offsetRectangle{src, i.x, i.y}, dest)
}

type offsetRectangle struct {
	Rectangle
	offsetX, offsetY int
}

func (r *offsetRectangle) TopLeft() (int, int) {
	x, y := r.Rectangle.TopLeft()
	return x + r.offsetX, y + r.offsetY
}
