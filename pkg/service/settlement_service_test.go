package service

import (
	"errors"
	"testing"

	"github.com/lackmus/settlementgengo/pkg/model"
)

type mockSettlementStorage struct {
	saved         []model.Settlement
	settlements   []model.Settlement
	failSave      bool
	failLoadAll   bool
	failDelete    bool
	failDeleteAll bool
}

func (m *mockSettlementStorage) LoadSettlement(name string) (model.Settlement, error) {
	for _, s := range m.settlements {
		if s.Name == name {
			return s, nil
		}
	}
	return model.Settlement{}, errors.New("not found")
}

func (m *mockSettlementStorage) LoadAllSettlements() ([]model.Settlement, error) {
	if m.failLoadAll {
		return nil, errors.New("load all failed")
	}
	return append([]model.Settlement(nil), m.settlements...), nil
}

func (m *mockSettlementStorage) SaveSettlement(settlement model.Settlement) error {
	if m.failSave {
		return errors.New("save failed")
	}
	m.saved = append(m.saved, settlement)
	return nil
}

func (m *mockSettlementStorage) SaveAllSettlements(settlements []model.Settlement) error {
	m.settlements = append([]model.Settlement(nil), settlements...)
	return nil
}

func (m *mockSettlementStorage) DeleteSettlement(name string) error {
	if m.failDelete {
		return errors.New("delete failed")
	}
	return nil
}

func (m *mockSettlementStorage) DeleteAllSettlements() error {
	if m.failDeleteAll {
		return errors.New("delete all failed")
	}
	m.settlements = nil
	return nil
}

func validSettlement(name string) model.Settlement {
	return model.Settlement{
		Name:       name,
		Faction:    "Marquise",
		XCoord:     10,
		YCoord:     20,
		Population: 150,
		Notes:      "Valid notes",
		NPCs:       []string{},
	}
}

func TestSettlementService_AddSettlementRejectsInvalid(t *testing.T) {
	storage := &mockSettlementStorage{}
	svc := SettlementService{Storage: storage, Settlements: []model.Settlement{}}

	invalid := validSettlement("Bad Input")
	invalid.XCoord = -10

	err := svc.AddSettlement(invalid)
	if err == nil {
		t.Fatal("AddSettlement() expected validation error, got nil")
	}
	if len(svc.Settlements) != 0 {
		t.Fatalf("AddSettlement() mutated in-memory list on error; got len=%d", len(svc.Settlements))
	}
	if len(storage.saved) != 0 {
		t.Fatalf("AddSettlement() saved invalid settlement; got save calls=%d", len(storage.saved))
	}
}

func TestSettlementService_AddSettlementSavesValid(t *testing.T) {
	storage := &mockSettlementStorage{}
	svc := SettlementService{Storage: storage, Settlements: []model.Settlement{}}

	settlement := validSettlement("Greenhall")
	err := svc.AddSettlement(settlement)
	if err != nil {
		t.Fatalf("AddSettlement() unexpected error: %v", err)
	}
	if len(svc.Settlements) != 1 {
		t.Fatalf("AddSettlement() expected 1 settlement in memory, got %d", len(svc.Settlements))
	}
	if len(storage.saved) != 1 {
		t.Fatalf("AddSettlement() expected 1 save call, got %d", len(storage.saved))
	}
}

func TestSettlementService_AddSettlement_SaveFailureDoesNotMutateInMemory(t *testing.T) {
	storage := &mockSettlementStorage{failSave: true}
	svc := SettlementService{Storage: storage, Settlements: []model.Settlement{}}

	err := svc.AddSettlement(validSettlement("Stonebridge"))
	if err == nil {
		t.Fatal("AddSettlement() expected save failure, got nil")
	}
	if len(svc.Settlements) != 0 {
		t.Fatalf("AddSettlement() mutated in-memory list on save failure; got len=%d", len(svc.Settlements))
	}
}

func TestSettlementService_UpdateSettlementRejectsInvalid(t *testing.T) {
	existing := validSettlement("Old Oak")
	storage := &mockSettlementStorage{}
	svc := SettlementService{Storage: storage, Settlements: []model.Settlement{existing}}

	updated := existing
	updated.Notes = "javascript:alert(1)"

	err := svc.UpdateSettlement(updated)
	if err == nil {
		t.Fatal("UpdateSettlement() expected validation error, got nil")
	}
	if svc.Settlements[0].Notes != existing.Notes {
		t.Fatalf("UpdateSettlement() should not mutate on validation failure; got notes=%q", svc.Settlements[0].Notes)
	}
	if len(storage.saved) != 0 {
		t.Fatalf("UpdateSettlement() saved invalid settlement; got save calls=%d", len(storage.saved))
	}
}

func TestSettlementService_UpdateSettlementNotFound(t *testing.T) {
	storage := &mockSettlementStorage{}
	svc := SettlementService{Storage: storage, Settlements: []model.Settlement{}}

	err := svc.UpdateSettlement(validSettlement("Missing"))
	if err == nil {
		t.Fatal("UpdateSettlement() expected not found error, got nil")
	}
}

func TestSettlementService_UpdateSettlement_SaveFailureDoesNotMutateInMemory(t *testing.T) {
	existing := validSettlement("Hillford")
	storage := &mockSettlementStorage{failSave: true}
	svc := SettlementService{Storage: storage, Settlements: []model.Settlement{existing}}

	updated := existing
	updated.Notes = "Updated"

	err := svc.UpdateSettlement(updated)
	if err == nil {
		t.Fatal("UpdateSettlement() expected save failure, got nil")
	}
	if svc.Settlements[0].Notes != existing.Notes {
		t.Fatalf("UpdateSettlement() mutated in-memory state on save failure; got notes=%q", svc.Settlements[0].Notes)
	}
}

func TestSettlementService_RemoveSettlement_DeleteFailureDoesNotMutateInMemory(t *testing.T) {
	existing := validSettlement("Riversend")
	storage := &mockSettlementStorage{failDelete: true}
	svc := SettlementService{Storage: storage, Settlements: []model.Settlement{existing}}

	err := svc.RemoveSettlement(existing.Name)
	if err == nil {
		t.Fatal("RemoveSettlement() expected delete failure, got nil")
	}
	if len(svc.Settlements) != 1 || svc.Settlements[0].Name != existing.Name {
		t.Fatalf("RemoveSettlement() mutated in-memory state on delete failure; got %+v", svc.Settlements)
	}
}

func TestSettlementService_DeleteAllSettlements_DeleteFailureDoesNotMutateInMemory(t *testing.T) {
	storage := &mockSettlementStorage{failDeleteAll: true}
	svc := SettlementService{Storage: storage, Settlements: []model.Settlement{validSettlement("Northwatch")}}

	err := svc.DeleteAllSettlements()
	if err == nil {
		t.Fatal("DeleteAllSettlements() expected delete all failure, got nil")
	}
	if len(svc.Settlements) != 1 {
		t.Fatalf("DeleteAllSettlements() mutated in-memory state on delete all failure; got len=%d", len(svc.Settlements))
	}
}

func TestSettlementService_GetSettlementsByFaction_FiltersResults(t *testing.T) {
	storage := &mockSettlementStorage{}
	svc := SettlementService{Storage: storage, Settlements: []model.Settlement{
		{Name: "A", Faction: "Marquise", XCoord: 1, YCoord: 1, Population: 100, Notes: "n", NPCs: []string{}},
		{Name: "B", Faction: "Eyrie", XCoord: 2, YCoord: 2, Population: 100, Notes: "n", NPCs: []string{}},
		{Name: "C", Faction: "Marquise", XCoord: 3, YCoord: 3, Population: 100, Notes: "n", NPCs: []string{}},
	}}

	results, err := svc.GetSettlementsByFaction("Marquise")
	if err != nil {
		t.Fatalf("GetSettlementsByFaction() unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("GetSettlementsByFaction() expected 2 settlements, got %d", len(results))
	}
	if results[0].Faction != "Marquise" || results[1].Faction != "Marquise" {
		t.Fatalf("GetSettlementsByFaction() returned wrong factions: %+v", results)
	}
}

func TestSettlementService_GetSettlementsByFaction_EmptyStateReturnsNoResults(t *testing.T) {
	storage := &mockSettlementStorage{}
	svc := SettlementService{Storage: storage}

	results, err := svc.GetSettlementsByFaction("Marquise")
	if err != nil {
		t.Fatalf("GetSettlementsByFaction() unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("GetSettlementsByFaction() expected 0 settlements, got %d", len(results))
	}
}
