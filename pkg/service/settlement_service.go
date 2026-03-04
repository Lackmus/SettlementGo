package service

import (
	"github.com/lackmus/settlementgengo/pkg/model"
	"github.com/lackmus/settlementgengo/pkg/shared"
)

type SettlementService struct {
	Settlements []model.Settlement
	Storage     shared.SettlementStorage
}

func NewSettlementService(storage shared.SettlementStorage) (*SettlementService, error) {
	settlements, err := storage.LoadAllSettlements()
	if err != nil {
		return nil, err
	}

	return &SettlementService{
		Settlements: settlements,
		Storage:     storage,
	}, nil
}

// add
func (s *SettlementService) AddSettlement(settlement model.Settlement) error {
	for i, existing := range s.Settlements {
		if existing.Name == settlement.Name {
			s.Settlements[i] = settlement
			return s.Storage.SaveSettlement(settlement)
		}
	}
	s.Settlements = append(s.Settlements, settlement)
	return s.Storage.SaveSettlement(settlement)
}

func (s *SettlementService) RemoveSettlement(name string) error {
	for i, settlement := range s.Settlements {
		if settlement.Name == name {
			s.Settlements = append(s.Settlements[:i], s.Settlements[i+1:]...)
			return s.Storage.DeleteSettlement(name)
		}
	}
	return nil
}

func (s *SettlementService) DeleteAllSettlements() error {
	s.Settlements = []model.Settlement{}
	return s.Storage.DeleteAllSettlements()
}

func (s *SettlementService) GetSettlement(name string) (model.Settlement, error) {
	return s.Storage.LoadSettlement(name)
}

func (s *SettlementService) GetAllSettlements() ([]model.Settlement, error) {
	return s.Storage.LoadAllSettlements()
}

func (s *SettlementService) UpdateSettlement(settlement model.Settlement) error {
	for i, settl := range s.Settlements {
		if settl.Name == settlement.Name {
			s.Settlements[i] = settlement
			return s.Storage.SaveSettlement(settlement)
		}
	}
	return nil
}
