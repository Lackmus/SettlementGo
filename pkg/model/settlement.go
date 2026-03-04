package model

// settlement struct
type Settlement struct {
	ID          int64    `json:"id"`
	Npcs        []string `json:"npcs"`
	Name        string   `json:"name"`
	Faction     string   `json:"faction"`
	XCoord      int      `json:"xCoord"`
	YCoord      int      `json:"yCoord"`
	Description string   `json:"description"`
	Population  int      `json:"population"`
}

func (s Settlement) AddNpc(n string) Settlement {
	s.Npcs = append(s.Npcs, n)
	return s
}

func (s Settlement) RemoveNpc(target string) Settlement {
	for i, n := range s.Npcs {
		if n == target {
			s.Npcs = append(s.Npcs[:i], s.Npcs[i+1:]...)
			break
		}
	}
	return s
}

func (s Settlement) PrintSettlement() {
	println("Settlement:", s.Name)
	println("Faction:", s.Faction)
	println("Coordinates:", s.XCoord, s.YCoord)
	println("NPCs:")
}
