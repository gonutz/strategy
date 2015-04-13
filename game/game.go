package game

func NewGame(imageLoader ImageLoader, camera SceneCamera) *Game {
	g := &Game{
		running: true,
	}
	g.state = NewPlayState(g, imageLoader, camera)
	return g
}

type Game struct {
	running   bool
	updating  bool
	state     GameState
	nextState GameState
}

func (g *Game) Running() bool {
	return g.running
}

func (g *Game) Quit() {
	g.running = false
}

func (g *Game) Update() {
	if g.state != nil {
		g.updating = true
		g.state.Update()
		g.updating = false
	}
	if g.nextState != nil {
		g.state = g.nextState
		g.nextState = nil
	}
}

func (g *Game) Draw() {
	if g.state != nil {
		g.state.Draw()
	}
}

func (g *Game) ChangeStateTo(state GameState) {
	g.nextState = state
}
