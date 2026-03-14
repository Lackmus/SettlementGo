package controllers

import (
	"errors"
	"strings"
	"testing"

	"github.com/lackmus/settlementgengo/pkg/model"
)

type mockNPCGateway struct {
	createNPCIDFn       func(npctype string, faction string) (string, error)
	createRandomNPCIDFn func() (string, error)
	deleteNPCByIDFn     func(id string) error
	getCreationOptions  func() CreationOptions
}

func (m *mockNPCGateway) CreateNPCAndID(npctype string, faction string) (string, error) {
	if m.createNPCIDFn == nil {
		return "", errors.New("CreateNPCAndID not configured")
	}
	return m.createNPCIDFn(npctype, faction)
}

func (m *mockNPCGateway) CreateRandomNPCAndID() (string, error) {
	if m.createRandomNPCIDFn == nil {
		return "", errors.New("CreateRandomNPCAndID not configured")
	}
	return m.createRandomNPCIDFn()
}

func (m *mockNPCGateway) DeleteNPC(id string) error {
	if m.deleteNPCByIDFn == nil {
		return nil
	}
	return m.deleteNPCByIDFn(id)
}

func (m *mockNPCGateway) GetCreationOptions() CreationOptions {
	if m.getCreationOptions == nil {
		return CreationOptions{}
	}
	return m.getCreationOptions()
}

func TestSettlementNPCProvider_GenerateRandomNPCsForSettlement_ValidatesInit(t *testing.T) {
	provider := newSettlementNPCProviderWithGateway(&mockNPCGateway{
		createRandomNPCIDFn: func() (string, error) {
			return "npc-1", nil
		},
		deleteNPCByIDFn: func(id string) error { return nil },
	})

	_, err := provider.GenerateRandomNPCsForSettlement(nil, 1)
	if err == nil {
		t.Fatal("GenerateRandomNPCsForSettlement() expected init error, got nil")
	}
	if !strings.Contains(err.Error(), "settlement is not initialized") {
		t.Fatalf("expected settlement init error, got: %v", err)
	}
}

func TestSettlementNPCProvider_GenerateRandomNPCsForSettlement_RollsBackOnFailure(t *testing.T) {
	settlement := model.Settlement{Name: "Oakwall", NPCs: []string{}}
	call := 0
	deleted := []string{}

	provider := newSettlementNPCProviderWithGateway(&mockNPCGateway{
		createRandomNPCIDFn: func() (string, error) {
			call++
			if call == 1 {
				return "npc-1", nil
			}
			return "", errors.New("generator failed")
		},
		deleteNPCByIDFn: func(id string) error {
			deleted = append(deleted, id)
			return nil
		},
	})

	_, err := provider.GenerateRandomNPCsForSettlement(&settlement, 2)
	if err == nil {
		t.Fatal("GenerateRandomNPCsForSettlement() expected generation failure, got nil")
	}
	if !strings.Contains(err.Error(), "generator failed") {
		t.Fatalf("expected generation failure, got: %v", err)
	}
	if len(deleted) != 1 || deleted[0] != "npc-1" {
		t.Fatalf("expected rollback delete for npc-1, got %v", deleted)
	}
	if len(settlement.NPCs) != 0 {
		t.Fatalf("expected settlement npc list rolled back to empty, got %v", settlement.NPCs)
	}
}

func TestSettlementNPCProvider_DeleteNPCFromSettlement_RemovesFromSettlement(t *testing.T) {
	settlement := model.Settlement{Name: "Thicket", NPCs: []string{"npc-a", "npc-b"}}
	provider := newSettlementNPCProviderWithGateway(&mockNPCGateway{})

	err := provider.DeleteNPCFromSettlement(&settlement, "npc-a")
	if err != nil {
		t.Fatalf("DeleteNPCFromSettlement() unexpected error: %v", err)
	}
	if len(settlement.NPCs) != 1 || settlement.NPCs[0] != "npc-b" {
		t.Fatalf("expected npc-a removed from settlement list, got %v", settlement.NPCs)
	}
}

func TestSettlementNPCProvider_DeleteNPC_ValidatesSettlement(t *testing.T) {
	provider := newSettlementNPCProviderWithGateway(&mockNPCGateway{})

	_, err := provider.DeleteNPC("npc-1", nil)
	if err == nil {
		t.Fatal("DeleteNPC() expected settlement init error, got nil")
	}
	if !strings.Contains(err.Error(), "settlement is not initialized") {
		t.Fatalf("expected settlement init error, got: %v", err)
	}
}

func TestSettlementNPCProvider_DeleteNPC_DoesNotMutateSettlementOnGatewayFailure(t *testing.T) {
	settlement := model.Settlement{Name: "Rook", NPCs: []string{"npc-a", "npc-b"}}
	provider := newSettlementNPCProviderWithGateway(&mockNPCGateway{
		deleteNPCByIDFn: func(id string) error {
			return errors.New("gateway delete failed")
		},
	})

	_, err := provider.DeleteNPC("npc-a", &settlement)
	if err == nil {
		t.Fatal("DeleteNPC() expected gateway delete failure, got nil")
	}
	if !strings.Contains(err.Error(), "gateway delete failed") {
		t.Fatalf("expected gateway delete failure, got: %v", err)
	}
	if len(settlement.NPCs) != 2 || settlement.NPCs[0] != "npc-a" || settlement.NPCs[1] != "npc-b" {
		t.Fatalf("expected settlement npc list unchanged on gateway failure, got %v", settlement.NPCs)
	}
}

func TestSettlementNPCProvider_GenerateRandomNPCsForSettlement_RollbackErrorJoinsCauseAndCleanup(t *testing.T) {
	settlement := model.Settlement{Name: "Oakwall", NPCs: []string{}}
	call := 0
	genErr := errors.New("generator failed")
	cleanupErr := errors.New("cleanup failed")

	provider := newSettlementNPCProviderWithGateway(&mockNPCGateway{
		createRandomNPCIDFn: func() (string, error) {
			call++
			if call == 1 {
				return "npc-1", nil
			}
			return "", genErr
		},
		deleteNPCByIDFn: func(id string) error {
			return cleanupErr
		},
	})

	_, err := provider.GenerateRandomNPCsForSettlement(&settlement, 2)
	if err == nil {
		t.Fatal("GenerateRandomNPCsForSettlement() expected failure, got nil")
	}
	if !errors.Is(err, genErr) {
		t.Fatalf("expected error to include generation cause, got: %v", err)
	}
	if !errors.Is(err, cleanupErr) {
		t.Fatalf("expected error to include rollback cleanup failure, got: %v", err)
	}
}

func TestSettlementNPCProvider_DeleteNPCBatch_RejectsEmptyID(t *testing.T) {
	deleteCalls := 0
	provider := newSettlementNPCProviderWithGateway(&mockNPCGateway{
		deleteNPCByIDFn: func(id string) error {
			deleteCalls++
			return nil
		},
	})

	err := provider.DeleteNPCBatch([]string{"", "npc-2"})
	if err == nil {
		t.Fatal("DeleteNPCBatch() expected empty id error, got nil")
	}
	if !strings.Contains(err.Error(), "npc id is empty at index 0") {
		t.Fatalf("expected indexed empty id error, got: %v", err)
	}
	if deleteCalls != 0 {
		t.Fatalf("expected no gateway delete calls after invalid input, got %d", deleteCalls)
	}
}
