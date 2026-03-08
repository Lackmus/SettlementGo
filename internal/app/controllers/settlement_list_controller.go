package controllers

import (
	"fmt"

	"github.com/lackmus/npcgengo"
	"github.com/lackmus/settlementgengo/pkg/model"
	"github.com/lackmus/settlementgengo/pkg/service"
	"github.com/lackmus/settlementgengo/pkg/shared"
)

type SettlementListController struct {
	SettlementService          service.SettlementService
	SettlementCreationSupplier service.SettlementCreationSupplier
	settlementNPCProvider      *SettlementNPCProvider
	observers                  []shared.SettlementObserver
}

func NewSettlementListController(
	settlementService service.SettlementService,
	settlementCreationsupplier service.SettlementCreationSupplier,
	npcGenerator npcgengo.NPCGen,
) *SettlementListController {
	settlementListController := &SettlementListController{
		SettlementService:          settlementService,
		SettlementCreationSupplier: settlementCreationsupplier,
		settlementNPCProvider:      NewSettlementNPCProvider(npcGenerator),
		observers:                  []shared.SettlementObserver{},
	}
	return settlementListController
}

func (c *SettlementListController) InitView(view shared.SettlementObserver) {
	c.RegisterObserver(view)
	c.NotifyObservers()
}

//**************************
// Observer pattern methods
//**************************

func (c *SettlementListController) RegisterObserver(observer shared.SettlementObserver) {
	c.observers = append(c.observers, observer)
}

func (c *SettlementListController) RemoveObserver(observer shared.SettlementObserver) {
	for i, obs := range c.observers {
		if obs == observer {
			c.observers = append(c.observers[:i], c.observers[i+1:]...)
			break
		}
	}
}

func (c *SettlementListController) NotifyObservers() {
	settlements := c.SettlementService.Settlements
	for _, observer := range c.observers {
		observer.Update(settlements)
	}
}

//**************************
// Settlement management methods
//**************************

func (c *SettlementListController) CreateSettlement(name string, faction string) (model.Settlement, error) {
	settlement := service.CreateSettlement(name, faction)
	return c.AddSettlement(settlement)
}

func (c *SettlementListController) CreateRandomSettlement() (model.Settlement, error) {
	settlement := service.CreateRandomSettlement(c.SettlementCreationSupplier)
	return c.AddSettlement(settlement)
}

func (c *SettlementListController) CreateRandomSettlementWithNPCs(npcCount int) (model.Settlement, error) {
	if npcCount < 0 {
		return model.Settlement{}, fmt.Errorf("npcCount cannot be negative")
	}

	settlement, err := c.CreateRandomSettlement()
	if err != nil {
		return model.Settlement{}, err
	}

	if npcCount == 0 {
		return settlement, nil
	}

	return c.AddRandomNPCsToSettlement(settlement.Name, npcCount)
}

func (c *SettlementListController) AddSettlement(settlement model.Settlement) (model.Settlement, error) {
	if err := c.SettlementService.AddSettlement(settlement); err != nil {
		return model.Settlement{}, err
	}
	c.NotifyObservers()
	return settlement, nil
}

func (c *SettlementListController) RemoveSettlement(name string) error {
	if err := c.SettlementService.RemoveSettlement(name); err != nil {
		return err
	}
	c.NotifyObservers()
	return nil
}

func (c *SettlementListController) RemoveAllSettlements() error {
	if err := c.SettlementService.DeleteAllSettlements(); err != nil {
		return err
	}
	c.NotifyObservers()
	return nil
}

func (c *SettlementListController) GetSettlement(name string) (model.Settlement, error) {
	return c.SettlementService.GetSettlement(name)
}

func (c *SettlementListController) GetAllSettlements() ([]model.Settlement, error) {
	return c.SettlementService.GetAllSettlements()
}

func (c *SettlementListController) UpdateSettlement(settlement model.Settlement) error {
	if err := c.SettlementService.UpdateSettlement(settlement); err != nil {
		return err
	}
	c.NotifyObservers()
	return nil
}

func (c *SettlementListController) GetSettlementsByFaction(faction string) ([]model.Settlement, error) {
	settlements, err := c.SettlementService.GetAllSettlements()
	if err != nil {
		return nil, err
	}
	var filtered []model.Settlement
	for _, s := range settlements {
		if s.Faction == faction {
			filtered = append(filtered, s)
		}
	}
	return filtered, nil
}

//**************************
// NPC management methods
//**************************

func (c *SettlementListController) AddNPCToSettlement(settlementName string, npctype string, faction string) (model.Settlement, error) {
	if c.settlementNPCProvider == nil {
		return model.Settlement{}, fmt.Errorf("npc generator is not configured")
	}
	settlement, err := c.GetSettlement(settlementName)
	if err != nil {
		return model.Settlement{}, err
	}

	settlement, err = c.settlementNPCProvider.GenerateNPCsForSettlement(&settlement, npctype, faction, 1)
	if err != nil {
		return model.Settlement{}, fmt.Errorf("failed to generate npc for settlement %q: %w", settlementName, err)
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

	originalNPCCount := len(settlement.Npcs)

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

func (c *SettlementListController) generatedNPCIDs(settlement model.Settlement, originalNPCCount int, expectedGeneratedCount int) ([]string, error) {
	generatedCount := len(settlement.Npcs) - originalNPCCount
	if generatedCount != expectedGeneratedCount {
		return nil, fmt.Errorf("expected %d generated NPCs, got %d", expectedGeneratedCount, generatedCount)
	}
	return append([]string(nil), settlement.Npcs[originalNPCCount:]...), nil
}

func (c *SettlementListController) rollbackGeneratedNPCs(generatedNPCIDs []string, cause error) error {
	if cleanupErr := c.settlementNPCProvider.DeleteNPCs(generatedNPCIDs); cleanupErr != nil {
		return fmt.Errorf("%w (rollback failed: %v)", cause, cleanupErr)
	}
	return cause
}

func (c *SettlementListController) GetNpcsInSettlement(name string) ([]string, error) {
	settlement, err := c.SettlementService.GetSettlement(name)
	if err != nil {
		return nil, err
	}
	return settlement.Npcs, nil
}
