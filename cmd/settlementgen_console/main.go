package main

import (
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

	viewer := console.NewConsoleView()
	controller := controllers.NewSettlementListController(*settlementService, viewer)

	controller.AddSettlement(
		service.CreateSettlement(
			"Test Settlement",
			"A test settlement for demonstration purposes.",
			100,
			"Test Faction",
			10.0, 20.0,
		),
	)

}
