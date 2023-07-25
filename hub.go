package main

import (
	"errors"
	"fmt"

	"github.com/h0rzn/sink-ships/game"
)

type Hub struct {
	// mu        *sync.Mutex
	games    map[*game.Game]map[string]*Client
	unhandled map[*Client]bool
	gameOver  chan string
}

func NewHub() *Hub {
	return &Hub{
		games: make(map[*game.Game]map[string]*Client),
		gameOver: make(chan string),
	}
}

func (h *Hub) CreateGame() *game.Game {
	game := game.NewGame()
	h.games[game] = make(map[string]*Client)
	return game
}

func (h *Hub) DeleteGame(gid string) (err error) {
	if game, exists := h.gameByID(gid); exists {
		delete(h.games, game)
	}
	return errors.New("hub -> delete game: cannot find game")
}

func (h *Hub) JoinGame(client *Client, gid string) (*game.Game, bool) {
	if game, exists := h.gameByID(gid); exists {
		h.games[game][client.ID] = client

		player, err := game.AddClient(client.ID)
		if err != nil {
			return nil, false
		}
		err = game.Join(player)
		if err != nil {
			return nil, false
		}
		return game, true
	}
	return nil, false
}

func (h *Hub) LeaveGame(gid string, client *Client) {
	if game, exists := h.gameByID(gid); exists {
		delete(h.games[game], client.ID)
	}
}

func (h *Hub) gameByID(id string) (*game.Game, bool) {
	for game := range h.games {
		if game.ID == id {
			return game, true
		}
	}
	return nil, false
}

func (h *Hub) Run() {
	for gid := range h.gameOver {
		// clean up game
		err := h.DeleteGame(gid)
		if err != nil {
			fmt.Println("err: ", err)
		}
	}
}

// func (h *Hub) randID() string {
// 	return ""
// }
