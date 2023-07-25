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
	Players      map[string]*Player
	ActivePlayer string
	State        int
	MovesIn      chan interface{}
}

// Server -> WS-Client
type Update struct {
	Cords    Cords
	Author   string
	Status   int
	GameOver bool
}

func NewGame() *Game {
	return &Game{
		ID:           "game-id",
		Players:      make(map[string]*Player),
		ActivePlayer: "123",
		State:        StateUnstarted,
		MovesIn:      make(chan interface{}),
	}
}

func (g *Game) Join(p *Player) error {
	if len(g.Players) >= 2 {
		return errors.New("max party size reached")
	}

	if _, exists := g.Players[p.ID]; !exists {
		return errors.New("duplicate player id")
	}

	g.Players[p.ID] = p
	return nil
}

func (g *Game) AddClient(id string) (*Player, error) {
	player := NewPlayer(id)
	return player, g.Join(player)

}

func (g *Game) Start() (err error) {
	if len(g.Players) != 2 {
		return errors.New("game requires 2 parties")
	}

	// send maps
	for id, player := range g.Players {
		enemy, err := g.otherPlayer(id)
		if err != nil {
			return errors.New("error messag here")
		}
		_, _ = enemy, player

		update := struct{}{}
		enemy.UpdateIn <- update
	}

	return
}

func (g *Game) End() {}

func (g *Game) Ready() bool {
	return false
}

func (g *Game) GameOver(update *Update) {
	// clean up
}

func (g *Game) shoot(cords Cords, player *Player) (update *Update, err error) {
	fmt.Printf("[game->shoot] %#v \n %#v\n", cords, player)

	if player.ID != g.ActivePlayer {
		return update, errors.New("moving-party mismatch")
	}

	fmt.Println("[game->shoot]", cords)
	enemy, err := g.otherPlayer(player.ID)
	if err != nil {
		return nil, errors.New("unkown enemy player")
	}
	cellStatus := player.Map.Shoot(cords)


	if cellStatus != CellRedundantShot {
		update := &Update{
			Cords:    cords,
			Author:   player.ID,
			Status:   cellStatus,
			GameOver: enemy.Map.ShipsSunken(),
		}
		return update, nil
	}
	return update, nil
}

func (g *Game) GetActivePlayerID() string {
	return g.ActivePlayer
}

func (g *Game) switchTurns() error {
	otherPlayer, err := g.otherPlayer(g.ActivePlayer)
	g.ActivePlayer = otherPlayer.ID
	return err
}

func (g *Game) otherPlayer(id string) (*Player, error) {
	for pid, p := range g.Players {
		if pid != id {
			return p, nil
		}
	}
	return nil, errors.New("failed to get other player")
}

func (g *Game) Run() {
	for move := range g.MovesIn {
		switch m := move.(type) {
		case PlaceShipsMove:
			if player, exists := g.Players[m.Author]; exists {
				player.Map.PlaceShips(m.Ships)

			}
		case ShootMove:
			if player, exists := g.Players[m.Author]; exists {
				update, err := g.shoot(*m.Cords, player)
				if err != nil {
					break
				}
				if enemy, err := g.otherPlayer(player.ID); err == nil {
					enemy.UpdateIn <- update

					// player sunk all enemies ships
					// game is over
					if update.GameOver {
						g.GameOver(update)
						return
					} else {
						g.switchTurns()
					}

				}
			}
		default:
		}
	}
}
