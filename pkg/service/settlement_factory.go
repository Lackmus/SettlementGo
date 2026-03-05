package service

import (
	"math/rand"

	"github.com/lackmus/settlementgengo/pkg/model"
)

const (
	DefaultFaction = "Default Faction"
	DefaultNotes   = "Default notes for the settlement."
	MinPopulation  = 100
	MaxPopulation  = 1000
	DefaultXCoord  = int(^uint(0) >> 1) // Max int value
	DefaultYCoord  = int(^uint(0) >> 1) // Max int value

)

func CreateSettlement(name string, faction string) model.Settlement {
	return model.Settlement{
		Name:       name,
		Population: MakeRandomPopulation(),
		Npcs:       []string{},
		Faction:    faction,
		XCoord:     DefaultXCoord,
		YCoord:     DefaultYCoord,
		Notes:      DefaultNotes,
	}
}

// create random settlement with random name and faction
func CreateRandomSettlement(settlementCreationSupplier SettlementCreationSupplier) model.Settlement {
	randomName := settlementCreationSupplier.GetRandomSettlementName()
	randomFaction := settlementCreationSupplier.GetRandomFaction()
	return model.Settlement{
		Name:       randomName,
		Population: MakeRandomPopulation(),
		Npcs:       []string{},
		Faction:    randomFaction,
		XCoord:     DefaultXCoord,
		YCoord:     DefaultYCoord,
		Notes:      DefaultNotes,
	}
}

func MakeRandomPopulation() int {
	return 100 + rand.Intn(900)
}
