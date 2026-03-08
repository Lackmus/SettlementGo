package controllers

import (
	"fmt"

	"github.com/lackmus/npcgengo"
)

// SettlementNPCGateway defines the minimal NPC operations needed by settlement orchestration.
type SettlementNPCGateway interface {
	CreateNPCAndID(npctype string, faction string) (string, error)
	CreateRandomNPCAndID() (string, error)
	DeleteNPC(id string) error
}

type settlementNPCGateway struct {
	npcGenerator npcgengo.NPCGen
}

func NewSettlementNPCGateway(npcGenerator npcgengo.NPCGen) SettlementNPCGateway {
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
		return nil
	}
	if g.npcGenerator.NPCListController == nil {
		return fmt.Errorf("npc generator is not configured")
	}
	g.npcGenerator.NPCListController.DeleteNPC(id)
	return nil
}
