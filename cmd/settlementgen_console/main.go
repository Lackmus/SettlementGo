package main

import (
	"fmt"

	"github.com/lackmus/npcgengo"
	"github.com/lackmus/settlementgengo/internal/app/controllers"
	"github.com/lackmus/settlementgengo/internal/platform/loaders"
	"github.com/lackmus/settlementgengo/pkg/service"
	"github.com/lackmus/settlementgengo/ui/console"
)

const (
	npcDataDir     = "./data"
	settlementDir  = "./data/settlement_database"
	settlementData = "./data/settlement_data"
)

func main() {
	settlementService, err := service.NewSettlementService(loaders.NewJSONSettlementStorage(settlementDir))
	if err != nil {
		panic(err)
	}

	npcGenerator, err := npcgengo.NewNPCGenWithDataDir(npcDataDir)
	if err != nil {
		panic(err)
	}

	npc, err := npcGenerator.NPCListController.CreateRandomNPCWithSeed(12345)
	if err != nil {
		panic(err)
	}

	viewer := console.NewConsoleView()

	loaders := loaders.NewJSONSettlementConfigLoader(settlementData)
	factions := npcGenerator.GetFactions()
	settlementCreationSupplier := service.NewSettlementCreationSupplier(loaders, factions)

	controller := controllers.NewSettlementListController(*settlementService, viewer, *settlementCreationSupplier)

	npc, err = npcGenerator.NPCListController.GetNPCByID("0")
	if err != nil {
		panic(err)
	}

	settlementA := service.CreateRandomSettlement(*settlementCreationSupplier)
	settlementB := service.CreateRandomSettlement(*settlementCreationSupplier)
	settlementA.AddNpc(npc.ID)
	settlementB.AddNpc(npc.ID)

	fmt.Println("Adding settlements...")
	controller.CreateRandomSettlement()
	controller.AddSettlement(settlementA)
	controller.AddSettlement(settlementB)

	fmt.Println("Deleting all settlements...")
	npcGenerator.NPCListController.DeleteAllNPCs()

	controller.RemoveAllSettlements()

}
