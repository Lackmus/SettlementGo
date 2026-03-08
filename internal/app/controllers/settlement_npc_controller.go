package controllers

import (
	"fmt"

	"github.com/lackmus/settlementgengo/pkg/model"
)

// NPC management methods are split into a dedicated file to keep
// settlement_list_controller.go focused on generic settlement operations.
func (c *SettlementListController) AddNPCToSettlement(settlementName string, npctype string, faction string) (model.Settlement, error) {
	if c.settlementNPCProvider == nil {
		return model.Settlement{}, fmt.Errorf("npc generator is not configured")
	}
	settlement, err := c.GetSettlement(settlementName)
	if err != nil {
		return model.Settlement{}, err
	}
	originalNPCCount := len(settlement.NPCs)
	settlement, err = c.settlementNPCProvider.GenerateSingleNPCForSettlement(&settlement, npctype, faction)
	if err != nil {
		return model.Settlement{}, fmt.Errorf("failed to generate npc for settlement %q: %w", settlementName, err)
	}
	generatedNPCIDs, err := c.generatedNPCIDs(settlement, originalNPCCount, 1)
	if err != nil {
		return model.Settlement{}, err
	}
	if err := c.UpdateSettlement(settlement); err != nil {
		return model.Settlement{}, c.rollbackGeneratedNPCs(generatedNPCIDs, fmt.Errorf("failed to update settlement after npc generation: %w", err))
	}
	return settlement, nil
}

func (c *SettlementListController) AddRandomNPCToSettlement(settlementName string) (model.Settlement, error) {
	if c.settlementNPCProvider == nil {
		return model.Settlement{}, fmt.Errorf("npc generator is not configured")
	}
	settlement, err := c.GetSettlement(settlementName)
	if err != nil {
		return model.Settlement{}, err
	}
	originalNPCCount := len(settlement.NPCs)
	settlement, err = c.settlementNPCProvider.GenerateSingleRandomNPCForSettlement(&settlement)
	if err != nil {
		return model.Settlement{}, fmt.Errorf("failed to generate random npc for settlement %q: %w", settlementName, err)
	}
	generatedNPCIDs, err := c.generatedNPCIDs(settlement, originalNPCCount, 1)
	if err != nil {
		return model.Settlement{}, err
	}
	if err := c.UpdateSettlement(settlement); err != nil {
		return model.Settlement{}, c.rollbackGeneratedNPCs(generatedNPCIDs, fmt.Errorf("failed to update settlement after random npc generation: %w", err))
	}
	return settlement, nil
}

func (c *SettlementListController) AddRandomNPCsToSettlement(name string, npcCount int) (model.Settlement, error) {
	if npcCount < 0 {
		return model.Settlement{}, fmt.Errorf("npcCount cannot be negative")
	}
	if c.settlementNPCProvider == nil {
		return model.Settlement{}, fmt.Errorf("npc generator is not configured")
	}
	settlement, err := c.GetSettlement(name)
	if err != nil {
		return model.Settlement{}, err
	}
	if npcCount == 0 {
		return settlement, nil
	}

	originalNPCCount := len(settlement.NPCs)

	settlement, err = c.settlementNPCProvider.GenerateRandomNPCsForSettlement(&settlement, npcCount)
	if err != nil {
		return model.Settlement{}, err
	}

	generatedNPCIDs, err := c.generatedNPCIDs(settlement, originalNPCCount, npcCount)
	if err != nil {
		return model.Settlement{}, err
	}

	if err := c.UpdateSettlement(settlement); err != nil {
		return model.Settlement{}, c.rollbackGeneratedNPCs(generatedNPCIDs, fmt.Errorf("failed to update settlement after npc generation: %w", err))
	}

	return settlement, nil
}

func (c *SettlementListController) GetNPCsInSettlement(name string) ([]string, error) {
	settlement, err := c.SettlementService.GetSettlement(name)
	if err != nil {
		return nil, err
	}
	return settlement.NPCs, nil
}

func (c *SettlementListController) generatedNPCIDs(settlement model.Settlement, originalNPCCount int, expectedGeneratedCount int) ([]string, error) {
	generatedCount := len(settlement.NPCs) - originalNPCCount
	if generatedCount != expectedGeneratedCount {
		return nil, fmt.Errorf("expected %d generated NPCs, got %d", expectedGeneratedCount, generatedCount)
	}
	return append([]string(nil), settlement.NPCs[originalNPCCount:]...), nil
}

func (c *SettlementListController) rollbackGeneratedNPCs(generatedNPCIDs []string, cause error) error {
	if cleanupErr := c.settlementNPCProvider.DeleteNPCBatch(generatedNPCIDs); cleanupErr != nil {
		return fmt.Errorf("%w (rollback failed: %v)", cause, cleanupErr)
	}
	return cause
}

func (c *SettlementListController) DeleteNPCFromSettlement(settlementName string, npcID string) error {
	if c.settlementNPCProvider == nil {
		return fmt.Errorf("npc generator is not configured")
	}
	settlement, err := c.GetSettlement(settlementName)
	if err != nil {
		return err
	}
	if err := c.settlementNPCProvider.DeleteNPCFromSettlement(&settlement, npcID); err != nil {
		return fmt.Errorf("failed to delete npc %q from settlement %q: %w", npcID, settlementName, err)
	}
	if err := c.UpdateSettlement(settlement); err != nil {
		return fmt.Errorf("failed to update settlement after npc deletion: %w", err)
	}
	return nil
}

func (c *SettlementListController) DeleteAllNPCsFromSettlement(settlementName string) error {
	if c.settlementNPCProvider == nil {
		return fmt.Errorf("npc generator is not configured")
	}
	settlement, err := c.GetSettlement(settlementName)
	if err != nil {
		return err
	}
	if err := c.settlementNPCProvider.DeleteAllNPCsFromSettlement(&settlement); err != nil {
		return fmt.Errorf("failed to delete npcs from settlement %q: %w", settlementName, err)
	}
	if err := c.UpdateSettlement(settlement); err != nil {
		return fmt.Errorf("failed to update settlement after npc deletion: %w", err)
	}
	return nil
}

func (c *SettlementListController) MoveNPCBetweenSettlements(sourceName string, targetName string, npcID string) error {
	if c.settlementNPCProvider == nil {
		return fmt.Errorf("npc generator is not configured")
	}
	sourceSettlement, err := c.GetSettlement(sourceName)
	if err != nil {
		return fmt.Errorf("failed to get source settlement %q: %w", sourceName, err)
	}
	targetSettlement, err := c.GetSettlement(targetName)
	if err != nil {
		return fmt.Errorf("failed to get target settlement %q: %w", targetName, err)
	}
	if err := c.settlementNPCProvider.DeleteNPCFromSettlement(&sourceSettlement, npcID); err != nil {
		return fmt.Errorf("failed to delete npc %q from source settlement %q: %w", npcID, sourceName, err)
	}
	if err := c.UpdateSettlement(sourceSettlement); err != nil {
		return fmt.Errorf("failed to update source settlement after npc deletion: %w", err)
	}
	if err := c.settlementNPCProvider.AddNPCToSettlement(&targetSettlement, npcID); err != nil {
		return fmt.Errorf("failed to add npc %q to target settlement %q: %w", npcID, targetName, err)
	}
	if err := c.UpdateSettlement(targetSettlement); err != nil {
		return fmt.Errorf("failed to update target settlement after npc addition: %w", err)
	}
	return nil
}
