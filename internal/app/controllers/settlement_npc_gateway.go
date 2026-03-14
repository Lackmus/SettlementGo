package controllers

import (
	"fmt"

	"github.com/lackmus/npcgengo"
	"github.com/lackmus/settlementgengo/internal/platform/helpers"
)

// SettlementNPCGateway defines the minimal NPC operations needed by settlement orchestration.
type SettlementNPCGateway interface {
	CreateNPCAndID(npctype string, faction string) (string, error)
	CreateRandomNPCAndID() (string, error)
	DeleteNPC(id string) error
	GetCreationOptions() CreationOptions
}

type settlementNPCGateway struct {
	npcGenerator npcgengo.NPCGen
}

type CreationOptions struct {
	Factions                []string            `json:"factions"`
	Species                 []string            `json:"species"`
	Traits                  []string            `json:"traits"`
	NpcTypes                []string            `json:"npcTypes"`
	NpcSubtypeForTypeMap    map[string][]string `json:"npcSubtypeForTypeMap"`
	NpcSpeciesForFactionMap map[string][]string `json:"npcSpeciesForFactionMap"`
}

func newSettlementNPCGateway(npcGenerator npcgengo.NPCGen) SettlementNPCGateway {
	return &settlementNPCGateway{npcGenerator: npcGenerator}
}

func (g *settlementNPCGateway) CreateNPCAndID(npctype string, faction string) (string, error) {
	if g.npcGenerator.NPCListController == nil {
		return "", fmt.Errorf("npc generator is not configured")
	}
	npc, err := g.npcGenerator.NPCListController.CreateNPC(npctype, faction)
	if err != nil {
		return "", err
	}
	if npc.ID == "" {
		return "", fmt.Errorf("generated npc id is empty")
	}
	return npc.ID, nil
}

func (g *settlementNPCGateway) CreateRandomNPCAndID() (string, error) {
	if g.npcGenerator.NPCListController == nil {
		return "", fmt.Errorf("npc generator is not configured")
	}
	npc, err := g.npcGenerator.NPCListController.CreateRandomNPC()
	if err != nil {
		return "", err
	}
	if npc.ID == "" {
		return "", fmt.Errorf("generated npc id is empty")
	}
	return npc.ID, nil
}

func (g *settlementNPCGateway) DeleteNPC(id string) error {
	if id == "" {
		return fmt.Errorf("npc id is empty")
	}
	if g.npcGenerator.NPCListController == nil {
		return fmt.Errorf("npc generator is not configured")
	}
	g.npcGenerator.NPCListController.DeleteNPC(id)
	return nil
}

func (g *settlementNPCGateway) GetCreationOptions() CreationOptions {
	if g.npcGenerator.NPCListController == nil {
		return CreationOptions{}
	}
	options := g.npcGenerator.NPCListController.GetCreationOptions()
	if options == nil {
		return CreationOptions{}
	}

	return CreationOptions{
		Factions:                append([]string(nil), options.Factions...),
		Species:                 append([]string(nil), options.Species...),
		Traits:                  append([]string(nil), options.Traits...),
		NpcTypes:                append([]string(nil), options.NpcTypes...),
		NpcSubtypeForTypeMap:    helpers.CopyStringSliceMap(options.NpcSubtypeForTypeMap),
		NpcSpeciesForFactionMap: helpers.CopyStringSliceMap(options.NpcSpeciesForFactionMap),
	}
}
