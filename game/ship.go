package game

type ShipSection struct {
	Cords Cords
	Hit   bool
}

type Ship struct {
	Type          string   `json:"type"`
	Range         [2]Cords `json:"range"`
	Sections      []*ShipSection
	SectionAmount int
	SectionHits   int
	Wreck         bool
}

type Ships []*Ship

func (s *Ship) Shoot(cords Cords) bool {
	for _, sec := range s.Sections {
		if sec.Cords.X == cords.X && sec.Cords.Y == cords.Y {
			return s.SetHit(sec)
		}
	}
	return false
}

func (s *Ship) SetHit(sec *ShipSection) (isWreck bool) {
	sec.Hit = true
	s.SectionHits++

	if s.SectionHits == s.SectionAmount {
		s.Wreck = true
		return true
	}
	return false
}

func (s *Ship) On(cords Cords) bool {
	return cords.X >= s.Range[0].X &&
		cords.X <= s.Range[1].X &&
		cords.Y >= s.Range[0].Y &&
		cords.Y <= s.Range[1].Y
}
