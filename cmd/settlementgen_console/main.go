package main

import (
	"fmt"

	"github.com/lackmus/settlementgengo/internal/app"
	"github.com/lackmus/settlementgengo/ui/console"
)

func main() {
	settlemntGenApp := app.NewSettlementGenApp()

	npcGenerator := settlemntGenApp.NpcGenerator
	controller := settlemntGenApp.SettlementController
	//settlementCreationSupplier := settlemntGenApp.SettlementCreationSupplier

	settlementViewer := console.NewConsoleView(controller)
	controller.InitView(settlementViewer)

	npc, err := npcGenerator.NPCListController.CreateRandomNPC()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created NPC 1: %s,%s", npc.Name(), npc.ID)

	/*
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
	*/npcGenerator.NPCListController.DeleteAllNPCs()
}
