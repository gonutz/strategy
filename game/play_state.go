package game

import (
	"github.com/mewmew/tmx"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"path/filepath"
)

func NewPlayState(context GameStateChanger, imageLoader ImageLoader, cam ScreenCamera) *PlayState {
	state := &PlayState{context: context, camera: newSceneCamera(cam)}
	state.loadMap(imageLoader, "crater_world.tmx")
	return state
}

type PlayState struct {
	context        GameStateChanger
	camera         SceneCamera
	tileMap        *tileMap
	mouseX, mouseY int
}

type GameStateChanger interface {
	ChangeStateTo(state GameState)
}

func (s *PlayState) Update() {
	x, y, w, h := s.camera.GetVisibleArea()
	xMargin, yMargin := s.getMoveDistance()
	worldX, worldY := s.camera.ScreenToWorld(s.mouseX, s.mouseY)
	dx, dy := 0, 0
	if worldX < x+xMargin {
		dx = -1
	}
	if worldX > x+w-xMargin {
		dx = 1
	}
	if worldY < y+yMargin {
		dy = -1
	}
	if worldY > y+h-yMargin {
		dy = 1
	}
	s.camera.Move(dx*xMargin, dy*xMargin)
}

func (s *PlayState) getMoveDistance() (dx, dy int) {
	_, _, w, h := s.camera.GetVisibleArea()
	dx = int(float32(w)*0.01 + 0.5)
	if dx == 0 {
		dx = 1
	}
	dy = int(float32(h)*0.01 + 0.5)
	if dy == 0 {
		dy = 1
	}
	return
}

func (s *PlayState) Draw() {
	s.tileMap.draw(s.camera.GetVisibleArea())
}

func (s *PlayState) KeyDown(key sdl.Keycode) {
	dx, dy := s.getMoveDistance()
	dx *= 3
	dy *= 3
	switch key {
	case sdl.K_LEFT:
		s.camera.Move(-dx, 0)
	case sdl.K_RIGHT:
		s.camera.Move(dx, 0)
	case sdl.K_DOWN:
		s.camera.Move(0, dy)
	case sdl.K_UP:
		s.camera.Move(0, -dy)
	}
}

func (s *PlayState) ScrollUp(x, y int) {
	s.camera.ZoomIn(x, y)
}

func (s *PlayState) ScrollDown(x, y int) {
	s.camera.ZoomOut(x, y)
}

func (s *PlayState) MouseMovedTo(x, y int) {
	s.mouseX, s.mouseY = x, y
}

func (s *PlayState) loadMap(imageLoader ImageLoader, ID string) {
	resourcePath := filepath.Join(
		os.Getenv("GOPATH"), "src", "github.com", "gonutz", "strategy", "rsc")
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
			gid := ground.GetGID(x, y)
			s.tileMap.tileAt(x, y).image = lookUpTileImage(tileImages, gid)
		}
	}

	mapW, mapH := s.tileMap.size()
	worldW, worldH := mapW*s.tileMap.tileW, mapH*s.tileMap.tileH
	s.camera.SetWorldSize(worldW, worldH)
	s.camera.SetFocus(worldW/2, worldH/2)
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
	tileX, tileY := m.clampTileCoordinates(left/m.tileW, top/m.tileH)
	rightTile, bottomTile := m.clampTileCoordinates(
		(left+width-1)/m.tileW,
		(top+height-1)/m.tileH,
	)

	for x := tileX; x <= rightTile; x++ {
		for y := tileY; y <= bottomTile; y++ {
			m.tileAt(x, y).image.Draw(
				rect{0, 0, m.tileW, m.tileH},
				rect{x*m.tileW - left, y*m.tileH - top, m.tileW, m.tileH},
			)
		}
	}
}

func (m *tileMap) clampTileCoordinates(x, y int) (int, int) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	w, h := m.size()
	if x >= w {
		x = w - 1
	}
	if y >= h {
		y = h - 1
	}
	return x, y
}

func (m *tileMap) tileAt(x, y int) *tile {
	return &m.tiles[x+y*m.width]
}

type tile struct {
	image Image
}

func newTileImage(imageLoader ImageLoader, path string,
	tileW, tileH, margin, spacing, firstID int) (*tileImage, error) {
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
