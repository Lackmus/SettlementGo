package mapper

import (
	"fmt"

	npcmapper "github.com/lackmus/npcgengo/pkg/mapper"
	h "github.com/lackmus/settlementgengo/internal/platform/helpers"
	"github.com/lackmus/settlementgengo/pkg/model"
)

type NPCResolver func(id string) (npcmapper.NPCInput, error)

type SettlementInputMapper struct {
	Name       string
	Notes      string
	Population int
	Faction    string
	XCoord     int
	YCoord     int
	NPCIDs     []string
}

type SettlementCreateInput struct {
	Name                  string `json:"name"`
	Faction               string `json:"faction"`
	XCoord                int    `json:"xCoord"`
	YCoord                int    `json:"yCoord"`
	Population            int    `json:"population"`
	Notes                 string `json:"notes"`
	InitialRandomNPCCount int    `json:"initialRandomNpcCount"`
}

type SettlementView struct {
	Name       string               `json:"name"`
	Faction    string               `json:"faction"`
	XCoord     int                  `json:"xCoord"`
	YCoord     int                  `json:"yCoord"`
	Population int                  `json:"population"`
	Notes      string               `json:"notes"`
	NPCs       []npcmapper.NPCInput `json:"npcs"`
}

func ToSettlementInput(settlementInput model.Settlement) SettlementInputMapper {
	return SettlementInputMapper{
		Name:       settlementInput.Name,
		Notes:      settlementInput.Notes,
		Population: settlementInput.Population,
		Faction:    settlementInput.Faction,
		XCoord:     settlementInput.XCoord,
		YCoord:     settlementInput.YCoord,
		NPCIDs:     settlementInput.NPCs,
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
		NPCs:       input.NPCIDs,
	}
}

func ToSettlementModelValidated(input SettlementInputMapper) (model.Settlement, error) {
	settlement := ToSettlementModel(input)
	if err := h.ValidateSettlement(settlement); err != nil {
		return model.Settlement{}, err
	}
	return settlement, nil
}

func ToSettlementView(input SettlementInputMapper, resolver NPCResolver) SettlementView {
	npcs := make([]npcmapper.NPCInput, 0, len(input.NPCIDs))

	for _, npcID := range input.NPCIDs {
		if resolver == nil {
			npcs = append(npcs, npcmapper.NPCInput{ID: npcID, Name: "NPC controller unavailable"})
			continue
		}

		npc, err := resolver(npcID)
		if err != nil {
			npcs = append(npcs, npcmapper.NPCInput{
				ID:    npcID,
				Name:  "Missing NPC",
				Notes: fmt.Sprintf("Failed to load NPC: %v", err),
			})
			continue
		}

		npcs = append(npcs, npc)
	}

	return SettlementView{
		Name:       input.Name,
		Faction:    input.Faction,
		XCoord:     input.XCoord,
		YCoord:     input.YCoord,
		Population: input.Population,
		Notes:      input.Notes,
		NPCs:       npcs,
	}
}
