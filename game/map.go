package game

import "fmt"

const (
	CellRaw = iota
	CellMiss
	CellHit
	CellFatalHit
	CellRedundantShot
)

type Cell struct {
	Cords     Cords
	Status    int
	isSection bool
}

type Map struct {
	Cells [][]*Cell
	Ships []*Ship
}

func NewMap() *Map {
	return &Map{
		Cells: make([][]*Cell, 4),
		Ships: make([]*Ship, 0),
	}
}

func (m *Map) Init() {
	fmt.Println("[map] initing")
	for x := 0; x < 4; x++ {
		m.Cells[x] = make([]*Cell, 4)
		for y := 0; y < 4; y++ {
			m.Cells[x][y] = &Cell{
				Cords: Cords{
					X: x,
					Y: y,
				},
				Status:    CellRaw,
				isSection: false,
			}
		}
	}
}

func (m *Map) PlaceShips(ships []Ship) {
	for _, ship := range ships {
		for _, section := range ship.Sections {
			secX := section.Cords.X
			secY := section.Cords.Y

			m.Cells[secX][secY].isSection = true
			m.Ships = append(m.Ships, &ship)
		}
	}
}

func (m *Map) GetShip(cords Cords) (*Ship, bool) {
	return &Ship{}, false
}

func (m *Map) Shoot(cords Cords) int {
	cell := m.Cells[cords.X][cords.Y]
	if cell.Status == CellHit {
		return CellRedundantShot
	}

	if ship, exists := m.GetShip(cords); exists {
		if ship.Shoot(cords) {
			return CellFatalHit
		} else {
			return CellHit
		}
	} else {
		cell.Status = CellMiss
		return CellMiss
	}
}

func (m *Map) ShipsSunken() bool {
	for _, ship := range m.Ships {
		if !ship.Wreck {
			return false
		}
	}
	return true
}
