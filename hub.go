package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/h0rzn/sink-ships/game"
)

type Hub struct {
	mu       *sync.Mutex
	games    map[*game.Game]map[string]*Client
	gameOver chan string
}

func NewHub() *Hub {
	return &Hub{
		games:    make(map[*game.Game]map[string]*Client),
		gameOver: make(chan string),
	}
}

func (h *Hub) CreateGame() *game.Game {
	h.mu.Lock()
	g := game.NewGame()
	h.games[g] = make(map[string]*Client)
	h.mu.Unlock()
	return g
}

func (h *Hub) DeleteGame(gid string) (err error) {
	h.mu.Lock()
	if game, exists := h.gameByID(gid); exists {
		delete(h.games, game)
	}
	h.mu.Unlock()
	return errors.New("hub -> delete game: cannot find game")
}

func (h *Hub) JoinGame(client *Client, gid string) (*game.Game, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if game, exists := h.gameByID(gid); exists {
		clientID := h.randID()
		h.games[game][clientID] = client

		player, err := game.AddClient(clientID)
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
	h.mu.Lock()
	if game, exists := h.gameByID(gid); exists {
		delete(h.games[game], client.ID)
	}
	h.mu.Unlock()
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

func (h *Hub) randID() string {
	return ""
}
