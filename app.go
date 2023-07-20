package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type App struct {
	Games   *GamePool
	Clients map[string]*Client
}

func NewApp() *App {
	return &App{
		Games:   NewGamePool(),
		Clients: make(map[string]*Client),
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

	fakeClientID := a.genID()
	c := NewClient(fakeClientID, con, a.Games)

	c.Read()
}

func (a *App) genID() string {
	return ""
}
