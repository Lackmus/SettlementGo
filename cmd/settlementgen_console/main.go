package main

import (
	"fmt"

	"github.com/lackmus/settlementgengo/internal/app"
	"github.com/lackmus/settlementgengo/ui/console"
)

func main() {
	settlementGenApp := app.NewSettlementGenApp()

	npcGenerator := settlementGenApp.NpcGenerator
	controller := settlementGenApp.SettlementController

	settlementViewer := console.NewConsoleView(controller)
	controller.InitView(settlementViewer)

	settlementA, err := controller.CreateRandomSettlementWithNPCs(2)
	if err != nil {
		panic(err)
	}
	settlementB, err := controller.CreateRandomSettlementWithNPCs(2)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created settlements with generated NPCs...")
	fmt.Printf("Created: %s, NPCs: %v\n", settlementA.Name, settlementA.NPCs)
	fmt.Printf("Created: %s, NPCs: %v\n", settlementB.Name, settlementB.NPCs)

	fmt.Println("Current settlements and their NPCs...")
	settlements := controller.SettlementService.Settlements
	for _, settlement := range settlements {
		fmt.Printf("Settlement: %s, NPCs: %v\n", settlement.Name, settlement.NPCs)
	}

	fmt.Println("Loading settlements and their NPCs...")
	for _, settlement := range settlements {
		fmt.Printf("Settlement: %s\n", settlement.Name)
		for _, npcID := range settlement.NPCs {
			npc, err := npcGenerator.NPCListController.GetNPCByID(npcID)
			if err != nil {
				fmt.Printf("  NPC ID: %s, Error: %v\n", npcID, err)
			} else {
				fmt.Printf("%s\n", npc.ShortString())
			}
		}
	}

	fmt.Println("Deleting all settlements...")

	if err := controller.RemoveAllSettlements(); err != nil {
		panic(err)
	}
	npcGenerator.NPCListController.DeleteAllNPCs()

	fmt.Println("Done.")
}
