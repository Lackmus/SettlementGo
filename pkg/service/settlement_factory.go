package service

import "github.com/lackmus/settlementgengo/pkg/model"

func CreateSettlement(name string, description string, population int, faction string, xCoord int, yCoord int) model.Settlement {
	return model.Settlement{
		Name:        name,
		Description: description,
		Population:  population,
		Faction:     faction,
		XCoord:      xCoord,
		YCoord:      yCoord,
		Npcs:        []string{},
	}
}
