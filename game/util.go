package game

type Cords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlaceShipsMove struct {
	Author string
	Ships []Ship
}

type ShootMove struct {
	Author string
	Cords *Cords
}