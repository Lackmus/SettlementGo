package service

import (
	"fmt"

	h "github.com/lackmus/settlementgengo/internal/platform/helpers"
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
	if err := h.ValidateSettlement(settlement); err != nil {
		return err
	}

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
	if err := h.ValidateSettlement(settlement); err != nil {
		return err
	}

	for i, settl := range s.Settlements {
		if settl.Name == settlement.Name {
			s.Settlements[i] = settlement
			return s.Storage.SaveSettlement(settlement)
		}
	}

	return fmt.Errorf("settlement %q not found", settlement.Name)
}
