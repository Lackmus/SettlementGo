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

func (c *SettlementListController) AddRandomNPCsToSettlement(name string, npcCount int) (model.Settlement, error) {
	if npcCount < 0 {
		return model.Settlement{}, fmt.Errorf("npcCount cannot be negative")
	}
	if c.settlementNPCProvider == nil || c.settlementNPCProvider.npcGenerator.NPCListController == nil {
		return model.Settlement{}, fmt.Errorf("npc generator is not configured")
	}
	settlement, err := c.GetSettlement(name)
	if err != nil {
		return model.Settlement{}, err
	}

	for i := 0; i < npcCount; i++ {
		settlement = *c.settlementNPCProvider.GenerateRandomNPCInSettlement(&settlement)
	}

	if err := c.UpdateSettlement(settlement); err != nil {
		return model.Settlement{}, err
	}

	return settlement, nil
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

func (c *SettlementListController) GetNpcsInSettlement(name string) ([]string, error) {
	settlement, err := c.SettlementService.GetSettlement(name)
	if err != nil {
		return nil, err
	}
	return settlement.Npcs, nil
}
