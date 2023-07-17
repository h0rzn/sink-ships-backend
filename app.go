package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/h0rzn/sink-ships/game"
	"github.com/h0rzn/sink-ships/netw"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type App struct {
	Games   map[string]*game.Game
	Clients map[string]*netw.Client
}

func NewApp() *App {
	return &App{
		Clients: make(map[string]*netw.Client),
	}
}

func (a *App) Run() {
	http.HandleFunc("/", a.wsHandler)
	http.ListenAndServe(":9090", nil)
}

func (a *App) wsHandler(w http.ResponseWriter, r *http.Request) {
	// defer close?
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("[app] handling ws")

	fakeRequestedGameID := "000"
	fakeClientID := "999"
	c := netw.NewClient(fakeClientID, con)
	if game, ok := a.Games[fakeRequestedGameID]; ok {
		c.Join(game)
		a.Clients[fakeClientID] = c
	}

}
