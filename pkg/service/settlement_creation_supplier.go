package service

import (
	"math/rand"

	"github.com/lackmus/settlementgengo/pkg/shared"
)

type SettlementCreationSupplier struct {
	settlementNames []string
	Factions        []string
}

func NewSettlementCreationSupplier(configLoader shared.SettlementConfigLoader, factions []string) *SettlementCreationSupplier {
	settlementNames, err := configLoader.LoadSettlementNames()
	if err != nil {
		panic(err)
	}
	return &SettlementCreationSupplier{
		settlementNames: settlementNames,
		Factions:        factions,
	}
}

// GetRandomSettlementName returns a random settlement name from the config loader
func (s *SettlementCreationSupplier) GetRandomSettlementName() string {
	if len(s.settlementNames) == 0 {
		return "Default Settlement Name"
	}
	return s.settlementNames[rand.Intn(len(s.settlementNames))]
}

func (s *SettlementCreationSupplier) GetRandomFaction() string {
	if len(s.Factions) == 0 {
		return "Default Faction"
	}
	return s.Factions[rand.Intn(len(s.Factions))]
}
