package main

import "encoding/json"

type RegisterFrame struct {
	Type string `json:"type"`
	Key string `json:"key"`
	Action string `json:"action"`
	PlayerName string `json:"player_name"`
	Data interface{}
}

type RegisterJoinData struct {
	GameID string `json:"game_id"`
}

type RegisterCreateData struct {
	
}

type RegisterResponse struct {
	Type string `json:"type"`
	PlayerID string `json:"id"`
	GameID string `json:"game_id"`
}

type ErrorMessage struct {}

type BaseMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type JoinData struct {
	GameID   string `json:"game_id"`
	Username string `json:"username"`
}

type CreateGameData struct {
	Name     string `json:"name"`
	Starting bool   `json:"self_starting"`
}

type PlaceData struct {
	Ships []PlaceShipsData
}

type PlaceShipsData struct {
	Type string `json:"ship_type"`
	From [2]int `json:"from"`
	To   [2]int `json:"to"`
}

type ShotData struct {
	X int
	Y int
}
