package anim

import "github.com/gonutz/strategy/images"

// TODO test all of this
// TODO what about sharing animations, e.g. several build animations that occur
// in different places and for different players, each will have its own trigger;
// creating the same animation over and over does not make sense, maybe solve
// this problem on a higher level, give the actual BuildAnimation an extra
// function for triggering creations, the question is then how to distringuish
// the multiple instances, maybe have the TileAnimation trigger have some kind
// of context parameter, so when the trigger is executed the info about who
// created it is still available?

type TileAnimation struct {
	frames        []tileFrame
	current       int
	timeInCurrent int
}

func NewTileAnimation(frameSets ...frameSet) *TileAnimation {
	return &TileAnimation{frames: Sequence(frameSets...).frames()}
}

func (a *TileAnimation) Update() {
	a.timeInCurrent++
	if a.timeInCurrent > a.frames[a.current].duration {
		a.nextFrame()
	}
}

func (a *TileAnimation) nextFrame() {
	if do := a.frames[a.current].trigger; do != nil {
		do()
	}
	a.current = (a.current + 1) % len(a.frames)
	a.timeInCurrent = 0
}

func (a *TileAnimation) DrawAt(x, y int) {
	if a.current < len(a.frames) {
		img := a.frames[a.current].image
		w, h := img.Size()
		img.Draw(
			images.Rect{0, 0, w, h},
			images.Rect{x, y, w, h},
		)
	}
}

type tileFrame struct {
	image    images.Image
	duration int
	trigger  func()
}

type frameSet interface {
	frames() []tileFrame
}

func TileRange(tiles *images.TiledImage, from, to int) frameSet {
	return &tileRange{tiles, from, to}
}

type tileRange struct {
	tiles    *images.TiledImage
	from, to int
}

func (r *tileRange) frames() []tileFrame {
	dx := 1
	if r.from > r.to {
		dx = -1
	}
	var frames []tileFrame
	for i := r.from; i != r.to; i += dx {
		frames = append(frames, tileFrame{r.tiles.Tile(i), 1, nil})
	}
	return frames
}

func Sequence(sets ...frameSet) frameSet {
	return &sequence{sets}
}

type sequence struct {
	sets []frameSet
}

func (s *sequence) frames() []tileFrame {
	var frames []tileFrame
	for _, set := range s.sets {
		frames = append(frames, set.frames()...)
	}
	return frames
}

func Loop(times int, set frameSet) frameSet {
	return &loop{times, set}
}

type loop struct {
	times int
	set   frameSet
}

func (l *loop) frames() []tileFrame {
	frames := l.set.frames()
	var looped []tileFrame
	for i := 0; i < l.times; i++ {
		looped = append(looped, frames...)
	}
	return looped
}

func After(set frameSet, do func()) frameSet {
	return &triggerAfter{set, do}
}

type triggerAfter struct {
	set frameSet
	do  func()
}

func (a *triggerAfter) frames() []tileFrame {
	frames := a.set.frames()
	if len(frames) > 0 {
		frames[len(frames)-1].trigger = a.do
	}
	return frames
}

func FrameDuration(d int, frameSets ...frameSet) frameSet {
	return &duration{d, frameSets}
}

type duration struct {
	time int
	sets []frameSet
}

func (d *duration) frames() []tileFrame {
	frames := Sequence(d.sets...).frames()
	for i := range frames {
		frames[i].duration = d.time
	}
	return frames
}
