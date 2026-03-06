package console

import (
	"fmt"

	"github.com/lackmus/settlementgengo/internal/app/controllers"
	"github.com/lackmus/settlementgengo/pkg/model"
)

type ConsoleView struct {
	controller  *controllers.SettlementListController
	settlements []model.Settlement
}

func NewConsoleView(ctrl *controllers.SettlementListController) *ConsoleView {
	return &ConsoleView{
		controller: ctrl,
	}
}

func (cv *ConsoleView) Update(settlements []model.Settlement) {
	cv.settlements = settlements
	cv.DisplaySettlements()
}

func (cv *ConsoleView) DisplaySettlements() {
	if len(cv.settlements) == 0 {
		fmt.Println("No settlements to display.")
		return
	}
	println("\nSettlements\n")
	for _, settlement := range cv.settlements {
		settlement.PrintSettlement()
		println("-----")
	}
}
