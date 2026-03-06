package app

import (
	"github.com/lackmus/npcgengo"
	"github.com/lackmus/settlementgengo/internal/app/controllers"
	"github.com/lackmus/settlementgengo/internal/platform/loaders"
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
	npcGenerator, err := npcgengo.NewNPCGenWithDataDir(defaultDataDir)
	if err != nil {
		panic(err)
	}
	settlemenService, err := service.NewSettlementService(loaders.NewJSONSettlementStorage(defaultDataDir + settlementDir))
	if err != nil {
		panic(err)
	}
	loaders := loaders.NewJSONSettlementConfigLoader(defaultDataDir + settlementData)
	factions := npcGenerator.GetFactions()
	settlementCreationSupplier := service.NewSettlementCreationSupplier(loaders, factions)
	settlementController := controllers.NewSettlementListController(*settlemenService, *settlementCreationSupplier)

	app := &SettlementGenApp{
		NpcGenerator:               *npcGenerator,
		SettlementService:          settlemenService,
		SettlementCreationSupplier: settlementCreationSupplier,
		SettlementController:       settlementController,
	}
	return app
}
