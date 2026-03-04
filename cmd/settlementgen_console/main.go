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
	settlementDir = "./data/settlement_database"
)

func main() {
	settlementService, err := service.NewSettlementService(loaders.NewJSONSettlementStorage(settlementDir))
	if err != nil {
		panic(err)
	}

	npcGenerator, err := npcgengo.NewNPCGenWithDataDir("./data")
	if err != nil {
		panic(err)
	}

	npc, err := npcGenerator.NPCListController.CreateRandomNPCWithSeed(12345)
	if err != nil {
		panic(err)
	}

	viewer := console.NewConsoleView()
	controller := controllers.NewSettlementListController(*settlementService, viewer)

	settlementA := service.CreateSettlement(
		"Test Settlement",
		"A test settlement for demonstration purposes.",
		100,
		"Test Faction",
		10, 20,
	)

	settlementB := service.CreateSettlement(
		"Another Settlement",
		"Another test settlement for demonstration purposes.",
		200,
		"Another Faction",
		30, 40,
	)

	npc, err = npcGenerator.NPCListController.GetNPCByID("0")
	if err != nil {
		panic(err)
	}

	fmt.Println("Adding NPC to settlements...")
	settlementA.AddNpc(npc.ID)
	settlementB.AddNpc(npc.ID)

	fmt.Println("Adding settlements...")
	fmt.Println("SettlementA:")
	controller.AddSettlement(settlementA)
	fmt.Println("SettlementB:")
	controller.AddSettlement(settlementB)

	fmt.Println("Deleting all settlements...")
	npcGenerator.NPCListController.DeleteAllNPCs()

	controller.RemoveAllSettlements()

}
