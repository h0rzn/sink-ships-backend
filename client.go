package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/h0rzn/sink-ships/game"
)

type Client struct {
	ID       string
	Con      *websocket.Conn
	Game     *game.Game
	Games    *GamePool
	UpdateIn chan interface{}
}

func NewClient(id string, con *websocket.Conn, games *GamePool) *Client {
	return &Client{
		ID:    id,
		Con:   con,
		Games: games,
	}
}

func (c *Client) Join(g *game.Game) (err error) {
	player := game.NewPlayer()
	c.Game = g
	if updates, err := c.Game.Join(player); err == nil {
		c.UpdateIn = updates

	}
	return
}

func (c *Client) Read() {
	defer c.Con.Close()

	// handle auth
	var authFrame *AuthFrame
	err := c.Con.ReadJSON(authFrame)
	if err != nil {
		// abort connection
		return
	}
	// validate auth

	for {
		var frame *BaseMessage
		err := c.Con.ReadJSON(&frame)
		if err != nil {
			fmt.Println("failed to read from websocket")
			return
		}
		fmt.Println("[client] out <- frame")

		c.HandleMessage(frame)

	}
}
func (c *Client) HandleMessage(frame *BaseMessage) (err error) {
	if c.Game.ActivePlayerID() != c.ID {
		return errors.New("ignoring move request: not my turn")
	}

	switch frame.Type {
	case "place":
		fmt.Println("[client] handling <place>")
		var data []game.Ship
		err = json.Unmarshal(frame.Data, &data)
		if err != nil {
			fmt.Println("[client] err", err)
			return
		}

		move := &game.PlaceShipsMove{
			Author: c.ID,
			Ships:  data,
		}
		c.Game.MovesIn <- move

	case "shoot":
		fmt.Println("[proxy] handling <shoot>")
		var data *game.Cords
		err = json.Unmarshal(frame.Data, &data)
		if err != nil {
			fmt.Println("[proxy] err", err)
			return err
		}

		move := &game.ShootMove{
			Author: c.ID,
			Cords:  data,
		}

		c.Game.MovesIn <- move

	default:
		fmt.Println("huhu")
	}

	return
}

func (c *Client) Send(msg interface{}) {
	fmt.Println("[client] sending update", msg)
}

func (c *Client) Close() {
	fmt.Println("[client] close by call")
}
