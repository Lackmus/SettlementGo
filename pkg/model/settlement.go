package model

import "fmt"

// settlement struct
type Settlement struct {
	NPCs       []string `json:"npcs"`
	Name       string   `json:"name"`
	Faction    string   `json:"faction"`
	XCoord     int      `json:"xCoord"`
	YCoord     int      `json:"yCoord"`
	Population int      `json:"population"`
	Notes      string   `json:"notes"`
}

func (s Settlement) HasNPC(id string) bool {
	for _, npcID := range s.NPCs {
		if npcID == id {
			return true
		}
	}
	return false
}

func (s *Settlement) AddNPC(npc string) error {
	if npc == "" {
		return fmt.Errorf("npc id cannot be empty")
	}
	if s.NPCs == nil {
		s.NPCs = []string{}
	}
	s.NPCs = append(s.NPCs, npc)
	return nil
}

func (s *Settlement) RemoveNPC(target string) {
	for i, n := range s.NPCs {
		if n == target {
			s.NPCs = append(s.NPCs[:i], s.NPCs[i+1:]...) //
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
	for _, npc := range s.NPCs {
		println(" -", npc)
	}
	println("Notes:", s.Notes)
}
