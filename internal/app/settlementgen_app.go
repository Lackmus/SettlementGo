package app

import (
	"github.com/lackmus/npcgengo"
	"github.com/lackmus/settlementgengo/internal/app/controllers"
	"github.com/lackmus/settlementgengo/internal/platform/loaders"
	"github.com/lackmus/settlementgengo/pkg/model"
	"github.com/lackmus/settlementgengo/pkg/service"
)

const (
	settlementDir  = "/settlement_database"
	settlementData = "/settlement_data"
	defaultDataDir = "./data"
)

type SettlementGenApp struct {
	NpcGenerator               npcgengo.NPCGen
	SettlementController       *controllers.SettlementListController
	SettlementCreationSupplier *service.SettlementCreationSupplier
	SettlementService          *service.SettlementService
}

func NewSettlementGenApp() *SettlementGenApp {
	return NewSettlementGenAppWithDataDir(defaultDataDir)
}

func NewSettlementGenAppWithDataDir(dir string) *SettlementGenApp {
	npcGenerator, err := npcgengo.NewNPCGenWithDataDir(dir)
	if err != nil {
		panic(err)
	}
	settlementService, err := service.NewSettlementService(loaders.NewJSONSettlementStorage(dir + settlementDir))
	if err != nil {
		panic(err)
	}
	loaders := loaders.NewJSONSettlementConfigLoader(dir + settlementData)
	factions := npcGenerator.GetFactions()
	settlementCreationSupplier := service.NewSettlementCreationSupplier(loaders, factions)
	settlementController := controllers.NewSettlementListController(*settlementService, *settlementCreationSupplier, *npcGenerator)

	app := &SettlementGenApp{
		NpcGenerator:               *npcGenerator,
		SettlementService:          settlementService,
		SettlementCreationSupplier: settlementCreationSupplier,
		SettlementController:       settlementController,
	}
	return app
}

// CreateRandomSettlementWithNPCs creates and saves a random settlement,
// then generates and attaches npcCount random NPC IDs to it.
func (a *SettlementGenApp) CreateRandomSettlementWithNPCs(npcCount int) (model.Settlement, error) {
	return a.SettlementController.CreateRandomSettlementWithNPCs(npcCount)
}

// AddRandomNPCsToSettlement generates npcCount random NPCs and appends their IDs
// to the named settlement. The updated settlement is persisted via controller update.
func (a *SettlementGenApp) AddRandomNPCsToSettlement(name string, npcCount int) (model.Settlement, error) {
	return a.SettlementController.AddRandomNPCsToSettlement(name, npcCount)
}
