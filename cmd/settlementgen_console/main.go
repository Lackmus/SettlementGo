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

	npc, err = npcGenerator.NPCListController.GetNPCByID("0")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Generated NPC ID %s:\n", npc.ID)

	settlementA.AddNpc(npc.ID)

	fmt.Printf("Settlement '%s' after adding NPC:\n", settlementA.Name)
	for _, npcID := range settlementA.Npcs {
		fmt.Printf("  NPC ID: %s\n", npcID)
	}

	controller.AddSettlement(settlementA)

	npcGenerator.NPCListController.DeleteAllNPCs()

	//controller.RemoveSettlement("Test Settlement")

}
