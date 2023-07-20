package main

import (
	"errors"
	"sync"

	"github.com/h0rzn/sink-ships/game"
)

type GamePool struct {
	mu    *sync.Mutex
	games map[string]*game.Game
}

func NewGamePool() *GamePool {
	return &GamePool{
		mu:    &sync.Mutex{},
		games: make(map[string]*game.Game),
	}
}

func (p *GamePool) Create() *game.Game {
	p.mu.Lock()

	g := &game.Game{}
	id := p.randID()
	p.games[id] = g

	p.mu.Unlock()
	return p.games[id]
}

func (p *GamePool) Remove(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	g, ok := p.games[id]
	if !ok {
		return errors.New("game remove: unkown game id")
	}
	// shutdown logic for game
	delete(p.games, g.ID)
	return nil
}

func (p *GamePool) randID() string {
	return ""
}
