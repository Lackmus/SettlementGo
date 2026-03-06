package main

import (
	"fmt"

	"github.com/lackmus/settlementgengo/internal/app"
	"github.com/lackmus/settlementgengo/pkg/service"
	"github.com/lackmus/settlementgengo/ui/console"
)

func main() {
	settlemntGenApp := app.NewSettlementGenApp()

	npcGenerator := settlemntGenApp.NpcGenerator
	controller := settlemntGenApp.SettlementController
	settlementCreationSupplier := settlemntGenApp.SettlementCreationSupplier

	settlementViewer := console.NewConsoleView(controller)
	controller.InitView(settlementViewer)

	npc, err := npcGenerator.NPCListController.CreateRandomNPC()
	if err != nil {
		panic(err)
	}

	settlementA := service.CreateRandomSettlement(*settlementCreationSupplier)
	settlementB := service.CreateRandomSettlement(*settlementCreationSupplier)
	settlementA.AddNpc(npc.ID)
	settlementB.AddNpc(npc.ID)

	fmt.Println("Adding settlements...")
	if err := controller.CreateRandomSettlement(); err != nil {
		panic(err)
	}
	if err := controller.AddSettlement(settlementA); err != nil {
		panic(err)
	}
	if err := controller.AddSettlement(settlementB); err != nil {
		panic(err)
	}

	fmt.Println("Deleting all settlements...")

	if err := controller.RemoveAllSettlements(); err != nil {
		panic(err)
	}
	npcGenerator.NPCListController.DeleteAllNPCs()
}
