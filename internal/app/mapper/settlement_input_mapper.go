package mapper

import (
	h "github.com/lackmus/settlementgengo/internal/platform/helpers"
	"github.com/lackmus/settlementgengo/pkg/model"
)

type SettlementInputMapper struct {
	Name       string
	Notes      string
	Population int
	Faction    string
	XCoord     int
	YCoord     int
	NPCIDs     []string
}

func ToSettlementInput(settlementInput model.Settlement) SettlementInputMapper {
	return SettlementInputMapper{
		Name:       settlementInput.Name,
		Notes:      settlementInput.Notes,
		Population: settlementInput.Population,
		Faction:    settlementInput.Faction,
		XCoord:     settlementInput.XCoord,
		YCoord:     settlementInput.YCoord,
		NPCIDs:     settlementInput.Npcs,
	}
}

func ToSettlementInputs(settlements []model.Settlement) []SettlementInputMapper {
	inputs := make([]SettlementInputMapper, len(settlements))
	for i, settlement := range settlements {
		inputs[i] = ToSettlementInput(settlement)
	}
	return inputs
}

func ToSettlementModel(input SettlementInputMapper) model.Settlement {
	return model.Settlement{
		Name:       input.Name,
		Notes:      input.Notes,
		Population: input.Population,
		Faction:    input.Faction,
		XCoord:     input.XCoord,
		YCoord:     input.YCoord,
		Npcs:       input.NPCIDs,
	}
}

func ToSettlementModelValidated(input SettlementInputMapper) (model.Settlement, error) {
	settlement := ToSettlementModel(input)
	if err := h.ValidateSettlement(settlement); err != nil {
		return model.Settlement{}, err
	}
	return settlement, nil
}
