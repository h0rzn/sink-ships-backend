package game

import (
	"errors"
	"fmt"
)

const (
	StateUnstarted = iota
	StatePlacing
	StatePlaying
	StateEnd
)

type Game struct {
	ID           string
	Players      map[*Player](chan interface{})
	ActivePlayer string
	State        int
	MovesIn      chan interface{}
}

type Update struct {
	Cords    Cords
	Author   string
	Status   int
	GameOver bool
}

func NewGame() *Game {
	return &Game{
		ID:           "game-id",
		Players:      make(map[*Player](chan interface{})),
		ActivePlayer: "123",
		State:        StateUnstarted,
		MovesIn:      make(chan interface{}),
	}
}

func (g *Game) Join(p *Player) (updates chan interface{}, err error) {
	if len(g.Players) >= 2 {
		return nil, errors.New("max party size reached")
	}

	for player := range g.Players {
		if player.ID == p.ID {
			return nil, errors.New("duplicate player id")
		}
	}

	g.Players[p] = make(chan interface{})
	return
}

func (g *Game) Start() (err error) {
	if len(g.Players) != 2 {
		return errors.New("game requires 2 parties")
	}

	return
}

func (g *Game) End() {}

func (g *Game) Ready() bool {
	return false
}

func (g *Game) shoot(cords Cords, player *Player) (update *Update, err error) {
	fmt.Printf("[game->shoot] %#v \n %#v\n", cords, player)

	if player.ID != g.ActivePlayer {
		return update, errors.New("moving-party mismatch")
	}

	fmt.Println("[game->shoot]", cords)
	cellStatus := player.Map.Shoot(cords)

	if cellStatus != CellRedundantShot {
		update := &Update{
			Cords:    cords,
			Author:   player.ID,
			Status:   cellStatus,
			GameOver: player.Map.ShipsSunken(),
		}
		return update, nil
	}
	return update, nil
}

func (g *Game) ActivePlayerID() string {
	return g.ActivePlayer
}

func (g *Game) GetWinner() bool {

	return true
}

func (g *Game) Run() {
	for move := range g.MovesIn {
		switch m := move.(type) {
		case nil:
			_ = m
		}
	}
}
