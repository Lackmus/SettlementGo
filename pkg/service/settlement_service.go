package service

import (
	"fmt"

	validation "github.com/lackmus/settlementgengo/internal/platform/helpers"
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

func (s *SettlementService) AddSettlement(settlement model.Settlement) error {
	if err := validation.ValidateSettlement(settlement); err != nil {
		return err
	}

	for i, existing := range s.Settlements {
		if existing.Name == settlement.Name {
			if err := s.Storage.SaveSettlement(settlement); err != nil {
				return err
			}
			s.Settlements[i] = settlement
			return nil
		}
	}
	if err := s.Storage.SaveSettlement(settlement); err != nil {
		return err
	}
	s.Settlements = append(s.Settlements, settlement)
	return nil
}

func (s *SettlementService) RemoveSettlement(name string) error {
	for i, settlement := range s.Settlements {
		if settlement.Name == name {
			if err := s.Storage.DeleteSettlement(name); err != nil {
				return err
			}
			s.Settlements = append(s.Settlements[:i], s.Settlements[i+1:]...)
			return nil
		}
	}
	return nil
}

func (s *SettlementService) DeleteAllSettlements() error {
	if err := s.Storage.DeleteAllSettlements(); err != nil {
		return err
	}
	s.Settlements = []model.Settlement{}
	return nil
}

func (s *SettlementService) GetSettlement(name string) (model.Settlement, error) {
	for _, settlement := range s.Settlements {
		if settlement.Name == name {
			return settlement, nil
		}
	}
	return model.Settlement{}, fmt.Errorf("settlement %q not found", name)
}

func (s *SettlementService) GetAllSettlements() ([]model.Settlement, error) {
	return append([]model.Settlement(nil), s.Settlements...), nil
}

func (s *SettlementService) GetSettlementsByFaction(faction string) ([]model.Settlement, error) {
	settlements, err := s.GetAllSettlements()
	if err != nil {
		return nil, err
	}

	filtered := make([]model.Settlement, 0, len(settlements))
	for _, settlement := range settlements {
		if settlement.Faction == faction {
			filtered = append(filtered, settlement)
		}
	}

	return filtered, nil
}

func (s *SettlementService) UpdateSettlement(settlement model.Settlement) error {
	if err := validation.ValidateSettlement(settlement); err != nil {
		return err
	}

	for i, currentSettlement := range s.Settlements {
		if currentSettlement.Name == settlement.Name {
			if err := s.Storage.SaveSettlement(settlement); err != nil {
				return err
			}
			s.Settlements[i] = settlement
			return nil
		}
	}

	return fmt.Errorf("settlement %q not found", settlement.Name)
}
