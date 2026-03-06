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
	controller.CreateRandomSettlement()
	controller.AddSettlement(settlementA)
	controller.AddSettlement(settlementB)

	fmt.Println("Deleting all settlements...")

	controller.RemoveAllSettlements()
	npcGenerator.NPCListController.DeleteAllNPCs()
}
