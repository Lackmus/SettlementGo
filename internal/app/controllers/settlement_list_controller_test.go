package controllers

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/lackmus/npcgengo"
	"github.com/lackmus/settlementgengo/pkg/model"
	"github.com/lackmus/settlementgengo/pkg/service"
)

func copyDir(t *testing.T, src, dst string) {
	t.Helper()

	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatalf("failed to read dir %s: %v", src, err)
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatalf("failed to create dir %s: %v", dst, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			copyDir(t, srcPath, dstPath)
			continue
		}

		srcFile, err := os.Open(srcPath)
		if err != nil {
			t.Fatalf("failed to open %s: %v", srcPath, err)
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			srcFile.Close()
			t.Fatalf("failed to create %s: %v", dstPath, err)
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			srcFile.Close()
			dstFile.Close()
			t.Fatalf("failed to copy %s to %s: %v", srcPath, dstPath, err)
		}

		srcFile.Close()
		dstFile.Close()
	}
}

func newTestNPCGen(t *testing.T) *npcgengo.NPCGen {
	t.Helper()

	dataDir := t.TempDir()
	copyDir(t, filepath.Clean("../../../data/creation_data"), filepath.Join(dataDir, "creation_data"))
	if err := os.MkdirAll(filepath.Join(dataDir, "npc_database"), 0o755); err != nil {
		t.Fatalf("failed to create npc database dir: %v", err)
	}

	npcGen, err := npcgengo.NewNPCGenWithDataDir(dataDir)
	if err != nil {
		t.Fatalf("failed to initialize npc generator: %v", err)
	}

	t.Cleanup(func() {
		npcGen.NPCListController.DeleteAllNPCs()
	})

	return npcGen
}

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
	for _, s := range m.saved {
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
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{}, npcgengo.NPCGen{})
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
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{}, npcgengo.NPCGen{})
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
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{}, npcgengo.NPCGen{})
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
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{}, npcgengo.NPCGen{})
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
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{}, npcgengo.NPCGen{})
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

func TestSettlementListController_CreateRandomSettlementWithNPCs_AppendsGeneratedNPCIDs(t *testing.T) {
	storage := &controllerMockStorage{all: []model.Settlement{}}
	svc := service.SettlementService{Storage: storage, Settlements: []model.Settlement{}}
	npcGen := newTestNPCGen(t)

	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{}, *npcGen)

	settlement, err := ctrl.CreateRandomSettlementWithNPCs(2)
	if err != nil {
		t.Fatalf("CreateRandomSettlementWithNPCs() unexpected error: %v", err)
	}
	if len(settlement.Npcs) != 2 {
		t.Fatalf("expected 2 NPC IDs, got %d", len(settlement.Npcs))
	}
	if settlement.Npcs[0] == "" || settlement.Npcs[1] == "" {
		t.Fatalf("expected non-empty NPC IDs, got %v", settlement.Npcs)
	}
}

func TestSettlementListController_AddRandomNPCsToSettlement_RequiresConfiguredNPCGenerator(t *testing.T) {
	existing := validControllerSettlement("Denwatch")
	storage := &controllerMockStorage{all: []model.Settlement{existing}}
	svc := service.SettlementService{Storage: storage, Settlements: []model.Settlement{existing}}
	ctrl := NewSettlementListController(svc, service.SettlementCreationSupplier{}, npcgengo.NPCGen{})

	_, err := ctrl.AddRandomNPCsToSettlement(existing.Name, 1)
	if err == nil {
		t.Fatal("AddRandomNPCsToSettlement() expected npc generator configuration error, got nil")
	}
}
