package controllers

import (
	"github.com/lackmus/settlementgengo/pkg/model"
	"github.com/lackmus/settlementgengo/pkg/service"
	"github.com/lackmus/settlementgengo/pkg/shared"
)

type SettlementListController struct {
	SettlementService service.SettlementService
	SettlementViewer  shared.SettlementViewer
	observers         []shared.SettlementObserver
}

func NewSettlementListController(service service.SettlementService, viewer shared.SettlementViewer) *SettlementListController {

	settlementListController := &SettlementListController{
		SettlementService: service,
		SettlementViewer:  viewer,
		observers:         []shared.SettlementObserver{},
	}
	settlementListController.RegisterObserver(viewer)
	return settlementListController

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

func (c *SettlementListController) AddSettlement(settlement model.Settlement) {
	c.SettlementService.AddSettlement(settlement)
	c.NotifyObservers()
}

func (c *SettlementListController) RemoveSettlement(name string) {
	c.SettlementService.RemoveSettlement(name)
	c.NotifyObservers()
}

func (c *SettlementListController) GetSettlement(name string) (model.Settlement, error) {
	return c.SettlementService.GetSettlement(name)
}

func (c *SettlementListController) GetAllSettlements() ([]model.Settlement, error) {
	return c.SettlementService.GetAllSettlements()
}

func (c *SettlementListController) UpdateSettlement(settlement model.Settlement) {
	c.SettlementService.UpdateSettlement(settlement)
	c.NotifyObservers()
}
