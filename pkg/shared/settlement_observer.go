package shared

import "github.com/lackmus/settlementgengo/pkg/model"

type SettlementObserver interface {
	Update(settlements []model.Settlement)
}
