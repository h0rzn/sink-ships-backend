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
	Hub    *Hub
	UpdateIn chan interface{}
}

func NewClient(id string, con *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:    id,
		Con:   con,
		Hub: hub,
	}
}

func (c *Client) Join(g *game.Game) (err error) {
	player := game.NewPlayer(c.ID)
	c.Game = g
	err = c.Game.Join(player)
	if err != nil {
		return err
	}

	return
}

func (c *Client) Read() {
	defer c.Con.Close()

	// handle register
	err := c.HandleRegister()
	if err != nil {
		// send error message
		// close client
		return
	}

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

func (c *Client) HandleRegister() error {
	var registerFrame *RegisterFrame
	err := c.Con.ReadJSON(registerFrame)
	if err != nil {
		return err
	}
	if registerFrame.Type != "register" {
		fmt.Println("first message is not of type 'register'")
	}
	if registerFrame.Key != "some_key" {
		fmt.Println("first message contains invalid register key")
	}
	switch registerFrame.Action {
	case "create":
		g := c.Hub.CreateGame()
		player, err := g.AddClient(c.ID)
		if err != nil {
			return err
		}
		_ = player
		c.Game = g
		// send game info to client
	case "join":
		var data *RegisterJoinData
		err := c.Con.ReadJSON(data)
		if err != nil {
			return err
		}

		if g, exists := c.Hub.JoinGame(c, data.GameID); exists {
			c.Game = g
			// send game info to client
			response := &RegisterResponse{
				Type: "register_response",
				PlayerID: c.ID,
				GameID: g.ID,
			}
			c.Send(response)
		} else {
			return err
		}

	default:
		return err
	}
	return nil
}

func (c *Client) HandleMessage(frame *BaseMessage) (err error) {
	if c.Game.GetActivePlayerID() != c.ID {
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
	}
	return
}

func (c *Client) Send(msg interface{}) {
	err := c.Con.WriteJSON(msg)
	if err != nil {
		// handle error
		return
	}
	fmt.Println("[client] send:", msg)
}

func (c *Client) Close() {
	fmt.Println("[client] close by call")
}
