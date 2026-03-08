package controllers

import (
	"fmt"

	"github.com/lackmus/npcgengo"
	"github.com/lackmus/settlementgengo/pkg/model"
)

type SettlementNPCProvider struct {
	gateway SettlementNPCGateway
}

func NewSettlementNPCProvider(npcGenerator npcgengo.NPCGen) *SettlementNPCProvider {
	return NewSettlementNPCProviderWithGateway(NewSettlementNPCGateway(npcGenerator))
}

func NewSettlementNPCProviderWithGateway(gateway SettlementNPCGateway) *SettlementNPCProvider {
	return &SettlementNPCProvider{gateway: gateway}
}

func (s *SettlementNPCProvider) AddNPCToSettlement(settlement *model.Settlement, npcID string) error {
	if err := s.validateProviderInput(settlement); err != nil {
		return err
	}
	return s.addNPCIDToSettlement(settlement, npcID)
}

func (s *SettlementNPCProvider) addNPCIDToSettlement(settlement *model.Settlement, npcID string) error {
	if npcID == "" {
		return nil
	}
	return settlement.AddNPC(npcID)
}

func (s *SettlementNPCProvider) DeleteNPCFromSettlement(settlement *model.Settlement, npcID string) error {
	if err := s.validateProviderInput(settlement); err != nil {
		return err
	}
	if npcID == "" {
		return nil
	}
	settlement.RemoveNPC(npcID)
	return nil
}

func (s *SettlementNPCProvider) DeleteAllNPCsFromSettlement(settlement *model.Settlement) error {
	if err := s.validateProviderInput(settlement); err != nil {
		return err
	}
	for _, npcID := range append([]string(nil), settlement.NPCs...) {
		if err := s.DeleteNPCFromSettlement(settlement, npcID); err != nil {
			return fmt.Errorf("failed to delete npc %q from settlement: %w", npcID, err)
		}
	}
	return nil
}

func (s *SettlementNPCProvider) DeleteNPC(npcID string, settlement *model.Settlement) (*model.Settlement, error) {
	if s == nil {
		return &model.Settlement{}, fmt.Errorf("settlement npc provider is not initialized")
	}
	if s.gateway == nil {
		return &model.Settlement{}, fmt.Errorf("npc generator is not configured")
	}
	if npcID == "" {
		return &model.Settlement{}, fmt.Errorf("npc id is empty")
	}
	settlement.RemoveNPC(npcID)
	return settlement, s.gateway.DeleteNPC(npcID)
}

func (s *SettlementNPCProvider) DeleteNPCBatch(ids []string) error {
	if s == nil {
		return fmt.Errorf("settlement npc provider is not initialized")
	}
	if s.gateway == nil {
		return fmt.Errorf("npc generator is not configured")
	}
	for _, id := range ids {
		if id == "" {
			continue
		}
		if err := s.gateway.DeleteNPC(id); err != nil {
			return err
		}
	}
	return nil
}

func (s *SettlementNPCProvider) validateProviderInput(settlement *model.Settlement) error {
	if s == nil {
		return fmt.Errorf("settlement NPC provider is not initialized")
	}
	if s.gateway == nil {
		return fmt.Errorf("npc generator is not configured")
	}
	if settlement == nil {
		return fmt.Errorf("settlement is not initialized")
	}
	return nil
}

func (s *SettlementNPCProvider) GenerateSingleNPCForSettlement(settlement *model.Settlement, npctype string, faction string) (model.Settlement, error) {
	if err := s.validateProviderInput(settlement); err != nil {
		return model.Settlement{}, err
	}
	npcID, err := s.gateway.CreateNPCAndID(npctype, faction)
	if err != nil {
		return model.Settlement{}, fmt.Errorf("failed to generate npc for settlement %q: %w", settlement.Name, err)
	}
	if err := s.addNPCIDToSettlement(settlement, npcID); err != nil {
		return model.Settlement{}, fmt.Errorf("failed to add generated npc to settlement %q: %w", settlement.Name, err)
	}
	if npcID == "" {
		return model.Settlement{}, fmt.Errorf("generated npc id is empty")
	}
	return *settlement, nil
}

func (s *SettlementNPCProvider) GenerateSingleRandomNPCForSettlement(settlement *model.Settlement) (model.Settlement, error) {
	if err := s.validateProviderInput(settlement); err != nil {
		return model.Settlement{}, err
	}
	npcID, err := s.gateway.CreateRandomNPCAndID()
	if err != nil {
		return model.Settlement{}, fmt.Errorf("failed to generate npc for settlement %q: %w", settlement.Name, err)
	}
	if err := s.addNPCIDToSettlement(settlement, npcID); err != nil {
		return model.Settlement{}, fmt.Errorf("failed to add generated npc to settlement %q: %w", settlement.Name, err)
	}
	if npcID == "" {
		return model.Settlement{}, fmt.Errorf("generated npc id is empty")
	}
	return *settlement, nil
}

func (s *SettlementNPCProvider) GenerateNPCsForSettlement(settlement *model.Settlement, npctype string, faction string, count int) (model.Settlement, error) {
	if err := s.validateProviderInput(settlement); err != nil {
		return model.Settlement{}, err
	}
	generatedNPCIDs := []string{}
	for i := 0; i < count; i++ {
		npcID, err := s.gateway.CreateNPCAndID(npctype, faction)
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
		npcID, err := s.gateway.CreateRandomNPCAndID()
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
	if err := s.addNPCIDToSettlement(settlement, npcID); err != nil {
		return generatedNPCIDs, s.rollbackGeneratedNPCBatch(settlement, generatedNPCIDs, err)
	}
	generatedNPCIDs = append(generatedNPCIDs, npcID)
	return generatedNPCIDs, nil
}

func (s *SettlementNPCProvider) rollbackGeneratedNPCBatch(settlement *model.Settlement, generatedNPCIDs []string, cause error) error {
	if cleanupErr := s.DeleteNPCBatch(generatedNPCIDs); cleanupErr != nil {
		return fmt.Errorf("%w (rollback failed: %v)", cause, cleanupErr)
	}
	if settlement != nil {
		for _, id := range generatedNPCIDs {
			settlement.RemoveNPC(id)
		}
	}
	return cause
}
