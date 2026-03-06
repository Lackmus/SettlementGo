package controllers

import (
	"errors"
	"testing"

	"github.com/lackmus/settlementgengo/pkg/model"
	"github.com/lackmus/settlementgengo/pkg/service"
)

type controllerMockStorage struct {
	failSave      bool
	failDeleteAll bool
	saved         []model.Settlement
	all           []model.Settlement
}

func (m *controllerMockStorage) LoadSettlement(name string) (model.Settlement, error) {
	for _, s := range m.all {
		if s.Name == name {
			return s, nil
		}
	}
	return model.Settlement{}, errors.New("not found")
}

func (m *controllerMockStorage) LoadAllSettlements() ([]model.Settlement, error) {
	return append([]model.Settlement(nil), m.all...), nil
}

func (m *controllerMockStorage) SaveSettlement(settlement model.Settlement) error {
	if m.failSave {
		return errors.New("save failed")
	}
	m.saved = append(m.saved, settlement)
	return nil
}

func (m *controllerMockStorage) SaveAllSettlements(settlements []model.Settlement) error {
	m.all = append([]model.Settlement(nil), settlements...)
	return nil
}

func (m *controllerMockStorage) DeleteSettlement(name string) error {
	return nil
}

func (m *controllerMockStorage) DeleteAllSettlements() error {
	if m.failDeleteAll {
		return errors.New("delete all failed")
	}
	m.all = nil
	return nil
}

type mockObserver struct {
	updates int
	lastLen int
}

func (o *mockObserver) Update(settlements []model.Settlement) {
	o.updates++
	o.lastLen = len(settlements)
}

func validControllerSettlement(name string) model.Settlement {
	return model.Settlement{
		Name:       name,
		Faction:    "Marquise",
		XCoord:     10,
		YCoord:     20,
		Population: 150,
		Notes:      "Valid notes",
		Npcs:       []string{},
	}
}

func TestSettlementListController_AddSettlement_ReturnsValidationErrorAndDoesNotNotify(t *testing.T) {
	storage := &controllerMockStorage{}
	svc := service.SettlementService{Storage: storage, Settlements: []model.Settlement{}}
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{})
	obs := &mockObserver{}
	ctrl.RegisterObserver(obs)

	invalid := validControllerSettlement("Broken")
	invalid.XCoord = -1

	settlement, err := ctrl.AddSettlement(invalid)
	if err == nil {
		t.Fatal("AddSettlement() expected validation error, got nil")
	}
	if obs.updates != 0 {
		t.Fatalf("observer should not be notified on error; got updates=%d", obs.updates)
	}
	if settlement.Name != "" {
		t.Fatalf("AddSettlement() expected zero value on error, got %+v", settlement)
	}
}

func TestSettlementListController_AddSettlement_NotifiesOnSuccess(t *testing.T) {
	storage := &controllerMockStorage{}
	svc := service.SettlementService{Storage: storage, Settlements: []model.Settlement{}}
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{})
	obs := &mockObserver{}
	ctrl.RegisterObserver(obs)

	npc, err := ctrl.AddSettlement(validControllerSettlement("Greenglade"))
	if err != nil {
		t.Fatalf("AddSettlement() unexpected error: %v", err)
	}
	if obs.updates != 1 {
		t.Fatalf("observer should be notified once on success; got updates=%d", obs.updates)
	}
	if obs.lastLen != 1 {
		t.Fatalf("observer expected 1 settlement after add; got %d", obs.lastLen)
	}
	if npc.Name != "Greenglade" {
		t.Fatalf("AddSettlement() expected name 'Greenglade', got '%s'", npc.Name)
	}
}

func TestSettlementListController_RemoveAllSettlements_PropagatesErrorAndSkipsNotify(t *testing.T) {
	storage := &controllerMockStorage{failDeleteAll: true}
	svc := service.SettlementService{Storage: storage, Settlements: []model.Settlement{validControllerSettlement("A")}}
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{})
	obs := &mockObserver{}
	ctrl.RegisterObserver(obs)

	err := ctrl.RemoveAllSettlements()
	if err == nil {
		t.Fatal("RemoveAllSettlements() expected error, got nil")
	}
	if obs.updates != 0 {
		t.Fatalf("observer should not be notified on error; got updates=%d", obs.updates)
	}
}

func TestSettlementListController_UpdateSettlement_NotFoundReturnsErrorAndSkipsNotify(t *testing.T) {
	storage := &controllerMockStorage{}
	svc := service.SettlementService{Storage: storage, Settlements: []model.Settlement{}}
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{})
	obs := &mockObserver{}
	ctrl.RegisterObserver(obs)

	err := ctrl.UpdateSettlement(validControllerSettlement("Missing"))
	if err == nil {
		t.Fatal("UpdateSettlement() expected not-found error, got nil")
	}
	if obs.updates != 0 {
		t.Fatalf("observer should not be notified on error; got updates=%d", obs.updates)
	}
}

func TestSettlementListController_CreateRandomSettlement_NotifiesObserver(t *testing.T) {
	storage := &controllerMockStorage{}
	svc := service.SettlementService{Storage: storage, Settlements: []model.Settlement{}}
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{})
	obs := &mockObserver{}
	ctrl.RegisterObserver(obs)

	npc, err := ctrl.CreateRandomSettlement()
	if err != nil {
		t.Fatalf("CreateRandomSettlement() unexpected error: %v", err)
	}
	if obs.updates != 1 {
		t.Fatalf("observer should be notified once on success; got updates=%d", obs.updates)
	}
	if obs.lastLen != 1 {
		t.Fatalf("observer expected 1 settlement after create random; got %d", obs.lastLen)
	}
	if npc.Name == "" {
		t.Fatal("CreateRandomSettlement() expected non-empty name")
	}
}
