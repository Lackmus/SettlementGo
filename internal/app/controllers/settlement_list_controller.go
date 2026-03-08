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
	settlementCreationSupplier service.SettlementCreationSupplier,
	npcGenerator npcgengo.NPCGen,
) *SettlementListController {
	settlementListController := &SettlementListController{
		SettlementService:          settlementService,
		SettlementCreationSupplier: settlementCreationSupplier,
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
	if c.SettlementExists(settlement.Name) {
		return model.Settlement{}, fmt.Errorf("settlement with name %q already exists", settlement.Name)
	}
	return c.AddSettlement(settlement)
}

func (c *SettlementListController) CreateRandomSettlement() (model.Settlement, error) {
	settlement := service.CreateRandomSettlement(c.SettlementCreationSupplier)
	for c.SettlementExists(settlement.Name) {
		settlement = service.CreateRandomSettlement(c.SettlementCreationSupplier)
	}
	return c.AddSettlement(settlement)
}

// check if settlement already exists
func (c *SettlementListController) SettlementExists(name string) bool {
	_, err := c.GetSettlement(name)
	return err == nil
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
	return c.SettlementService.GetSettlementsByFaction(faction)
}
