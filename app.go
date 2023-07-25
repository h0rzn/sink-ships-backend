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
	Hub *Hub
}

func NewApp() *App {
	return &App{
		Hub: NewHub(),
	}
}

func (a *App) Run() {
	http.HandleFunc("/", a.wsHandler)
	http.ListenAndServe(":9090", nil)
}

func (a *App) wsHandler(w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer con.Close()

	fmt.Println("[app] handling ws")

	client := NewClient(con, a.Hub)
	go client.Read()
}
