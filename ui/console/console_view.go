package console

import (
	"fmt"

	"github.com/lackmus/settlementgengo/pkg/model"
)

type ConsoleView struct {
	settlements []model.Settlement
}

func NewConsoleView() *ConsoleView {
	return &ConsoleView{
		settlements: []model.Settlement{},
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
	for _, settlement := range cv.settlements {
		settlement.PrintSettlement()
		println("-----")
	}
}
