package main

import "errors"

type Party struct {
	ID     string
	Moves  int
	Points int
	Map    *Map
}

type Cords struct {
	X int
	Y int
}

type Update struct {

}

type Game struct {
	ID          string
	Parties     map[string]*Party
	ActiveParty string
}

func NewGame() *Game {
	return &Game{
		ID:      "game-id",
		Parties: map[string]*Party{},
	}
}

func (g *Game) RegisterParty(party *Party) (err error) {
	// select first registered party as a dirty implementation for now
	if len(g.Parties) == 0 {
		g.Parties[party.ID] = party
		return
	}

	if len(g.Parties) >= 2 {
		return errors.New("max party size reached")
	}

	if _, alreadyExists := g.Parties[party.ID]; alreadyExists {
		return errors.New("duplicate party id")
	}
	g.Parties[party.ID] = party


	return
}

func (g *Game) Start() (err error) {
	if len(g.Parties) != 2 {
		return errors.New("game requires 2 parties")
	}

	return
}

func (g *Game) Shoot(cords Cords, party Party) (Update, error) {
	if (party.ID != g.ActiveParty) {
		return Update{}, errors.New("moving-party mismatch")
	}
	
	cellStatus := party.Map.Shoot(cords)
	switch (cellStatus) {
	case CellMiss:
		g.NextMove()
		break;
	case CellHit:
		g.NextMove()
		break;
	case CellFatalHit:
		g.NextMove()
		break;
	case CellRedundantShot:
		g.NextMove()
		break;
	default:
		return Update{}, errors.New("unkown cell status after shot")
	}

	return Update{}, nil
}

// After move clean up
func (g *Game) NextMove() {
	// switch party


}