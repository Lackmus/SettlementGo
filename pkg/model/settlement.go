package model

import "fmt"

// settlement struct
type Settlement struct {
	Npcs       []string `json:"npcs"`
	Name       string   `json:"name"`
	Faction    string   `json:"faction"`
	XCoord     int      `json:"xCoord"`
	YCoord     int      `json:"yCoord"`
	Population int      `json:"population"`
	Notes      string   `json:"notes"`
}

func (s Settlement) HasNpc(id string) bool {
	for _, npcID := range s.Npcs {
		if npcID == id {
			return true
		}
	}
	return false
}

func (s *Settlement) AddNpc(npc string) error {
	if npc == "" {
		return fmt.Errorf("npc id cannot be empty")
	}
	if s.Npcs == nil {
		s.Npcs = []string{}
	}
	s.Npcs = append(s.Npcs, npc)
	return nil
}

func (s *Settlement) RemoveNpc(target string) {
	for i, n := range s.Npcs {
		if n == target {
			s.Npcs = append(s.Npcs[:i], s.Npcs[i+1:]...) //
			break
		}
	}
}

func (s Settlement) PrintSettlement() {
	println("Settlement:", s.Name)
	println("Faction:", s.Faction)
	println("Population:", s.Population)
	println("Coordinates:", s.XCoord, s.YCoord)
	println("NPCs:")
	for _, npc := range s.Npcs {
		println(" -", npc)
	}
	println("Notes:", s.Notes)
}
