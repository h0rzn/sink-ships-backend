package game

type Player struct {
	ID    string
	Moves int
	Map   *Map
	UpdateIn chan interface{}
}

func NewPlayer(id string) *Player {
	return &Player{
		ID:    "party-id-here",
		Moves: 0,
		Map:   NewMap(),
		UpdateIn: make(chan interface{}),
	}
}

func (p *Player) Init() {
	p.Map.Init()
}
