package shared

import (
	"github.com/lackmus/settlementgengo/pkg/model"
)

type SettlementStorage interface {
	LoadSettlement(name string) (model.Settlement, error)

	LoadAllSettlements() ([]model.Settlement, error)

	SaveSettlement(nsettlement model.Settlement) error

	SaveAllSettlements(settlements []model.Settlement) error

	DeleteSettlement(name string) error

	DeleteAllSettlements() error
}
