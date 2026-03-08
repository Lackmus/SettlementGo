package controllers

import (
	"fmt"

	"github.com/lackmus/npcgengo"
	"github.com/lackmus/settlementgengo/pkg/model"
)

type SettlementNPCProvider struct {
	createNPCFn       func(npctype string, faction string) (string, error)
	createRandomNPCFn func() (string, error)
	deleteNPCFn       func(id string) error
}

func NewSettlementNPCProvider(npcGenerator npcgengo.NPCGen) *SettlementNPCProvider {
	return &SettlementNPCProvider{
		createNPCFn: func(npctype string, faction string) (string, error) {
			if npcGenerator.NPCListController == nil {
				return "", fmt.Errorf("npc generator is not configured")
			}
			npc, err := npcGenerator.NPCListController.CreateNPC(npctype, faction)
			if err != nil {
				return "", err
			}
			if npc.ID == "" {
				return "", fmt.Errorf("generated npc id is empty")
			}
			return npc.ID, nil
		},
		createRandomNPCFn: func() (string, error) {
			if npcGenerator.NPCListController == nil {
				return "", fmt.Errorf("npc generator is not configured")
			}
			npc, err := npcGenerator.NPCListController.CreateRandomNPC()
			if err != nil {
				return "", err
			}
			if npc.ID == "" {
				return "", fmt.Errorf("generated npc id is empty")
			}
			return npc.ID, nil
		},
		deleteNPCFn: func(id string) error {
			if npcGenerator.NPCListController == nil {
				return fmt.Errorf("npc generator is not configured")
			}
			if id == "" {
				return nil
			}
			npcGenerator.NPCListController.DeleteNPC(id)
			return nil
		},
	}
}

func (s *SettlementNPCProvider) DeleteNPCFromSettlement(settlement *model.Settlement, npcID string) error {
	if err := s.validateProviderInput(settlement); err != nil {
		return err
	}
	if err := s.deleteNPCFn(npcID); err != nil {
		return err
	}
	settlement.RemoveNpc(npcID)
	return nil
}

func (s *SettlementNPCProvider) DeleteNPCs(ids []string) error {
	if s == nil {
		return fmt.Errorf("settlement npc provider is not initialized")
	}
	for _, id := range ids {
		if id == "" {
			continue
		}
		if err := s.deleteNPCFn(id); err != nil {
			return err
		}
	}
	return nil
}

func (s *SettlementNPCProvider) validateProviderInput(settlement *model.Settlement) error {
	if s == nil {
		return fmt.Errorf("settlement NPC provider is not initialized")
	}
	if settlement == nil {
		return fmt.Errorf("settlement is not initialized")
	}
	return nil
}

func (s *SettlementNPCProvider) GenerateNPCsForSettlement(settlement *model.Settlement, npctype string, faction string, count int) (model.Settlement, error) {
	if err := s.validateProviderInput(settlement); err != nil {
		return model.Settlement{}, err
	}
	generatedNPCIDs := []string{}
	for i := 0; i < count; i++ {
		npcID, err := s.createNPCFn(npctype, faction)
		generatedNPCIDs, err = s.appendGeneratedNPCIDOrRollback(settlement, generatedNPCIDs, npcID, err)
		if err != nil {
			return model.Settlement{}, fmt.Errorf("failed generating npc %d/%d: %w", i+1, count, err)
		}
	}
	return *settlement, nil
}

func (s *SettlementNPCProvider) GenerateRandomNPCsForSettlement(settlement *model.Settlement, count int) (model.Settlement, error) {
	if err := s.validateProviderInput(settlement); err != nil {
		return model.Settlement{}, err
	}
	generatedNPCIDs := []string{}
	for i := 0; i < count; i++ {
		npcID, err := s.createRandomNPCFn()
		generatedNPCIDs, err = s.appendGeneratedNPCIDOrRollback(settlement, generatedNPCIDs, npcID, err)
		if err != nil {
			return model.Settlement{}, fmt.Errorf("failed generating npc %d/%d: %w", i+1, count, err)
		}
	}
	return *settlement, nil
}

func (s *SettlementNPCProvider) appendGeneratedNPCIDOrRollback(settlement *model.Settlement, generatedNPCIDs []string, npcID string, err error) ([]string, error) {
	if err != nil {
		return generatedNPCIDs, s.rollbackGeneratedNPCBatch(settlement, generatedNPCIDs, err)
	}
	if npcID == "" {
		return generatedNPCIDs, s.rollbackGeneratedNPCBatch(settlement, generatedNPCIDs, fmt.Errorf("generated npc id is empty"))
	}
	if err := settlement.AddNpc(npcID); err != nil {
		return generatedNPCIDs, s.rollbackGeneratedNPCBatch(settlement, generatedNPCIDs, err)
	}
	generatedNPCIDs = append(generatedNPCIDs, npcID)
	return generatedNPCIDs, nil
}

func (s *SettlementNPCProvider) rollbackGeneratedNPCBatch(settlement *model.Settlement, generatedNPCIDs []string, cause error) error {
	if cleanupErr := s.DeleteNPCs(generatedNPCIDs); cleanupErr != nil {
		return fmt.Errorf("%w (rollback failed: %v)", cause, cleanupErr)
	}
	if settlement != nil {
		for _, id := range generatedNPCIDs {
			settlement.RemoveNpc(id)
		}
	}
	return cause
}
