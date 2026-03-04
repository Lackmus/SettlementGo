package model

import "fmt"

// settlement struct
type Settlement struct {
	Npcs        []string `json:"npcs"`
	Name        string   `json:"name"`
	Faction     string   `json:"faction"`
	XCoord      int      `json:"xCoord"`
	YCoord      int      `json:"yCoord"`
	Description string   `json:"description"`
	Population  int      `json:"population"`
}

func (s *Settlement) AddNpc(npc string) {
	if npc == "" {
		fmt.Println("Cannot add empty NPC to settlement.")
		return
	}
	if s.Npcs == nil {
		s.Npcs = []string{}
	}
	s.Npcs = append(s.Npcs, npc)
}

func (s *Settlement) RemoveNpc(target string) {
	for i, n := range s.Npcs {
		if n == target {
			s.Npcs = append(s.Npcs[:i], s.Npcs[i+1:]...)
			break
		}
	}
}

func (s Settlement) PrintSettlement() {
	println("Settlement:", s.Name)
	println("Faction:", s.Faction)
	println("Coordinates:", s.XCoord, s.YCoord)
	println("NPCs:")
	for _, npc := range s.Npcs {
		println(" -", npc)
	}
}
