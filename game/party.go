package game

type Player struct {
	ID    string
	Moves int
	Map   *Map
}

func NewPlayer() *Player {
	return &Player{
		ID:    "party-id-here",
		Moves: 0,
		Map:   NewMap(),
	}
}

func (p *Player) Init() {
	p.Map.Init()
}
