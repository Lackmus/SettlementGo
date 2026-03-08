package controllers

import (
	"github.com/lackmus/npcgengo"
	"github.com/lackmus/settlementgengo/pkg/model"
)

type SettlementNPCProvider struct {
	npcGenerator npcgengo.NPCGen
}

func NewSettlementNPCProvider(npcGenerator npcgengo.NPCGen) *SettlementNPCProvider {
	return &SettlementNPCProvider{
		npcGenerator: npcGenerator,
	}
}

func (s *SettlementNPCProvider) GenerateNPCInSettlement(settlement *model.Settlement, npctype string, faction string) *model.Settlement {
	npc, err := s.npcGenerator.NPCListController.CreateNPC(npctype, faction)
	if err != nil {
		return settlement
	}
	settlement.AddNpc(npc.ID)
	return settlement
}

func (s *SettlementNPCProvider) GenerateRandomNPCInSettlement(settlement *model.Settlement) *model.Settlement {
	npc, err := s.npcGenerator.NPCListController.CreateRandomNPC()
	if err != nil {
		return settlement
	}
	settlement.AddNpc(npc.ID)
	return settlement
}
