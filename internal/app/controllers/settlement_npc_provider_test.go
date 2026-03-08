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

func TestSettlementNPCProvider_GenerateRandomNPCsForSettlement_ValidatesInit(t *testing.T) {
	provider := NewSettlementNPCProviderWithGateway(&mockNPCGateway{
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
	settlement := model.Settlement{Name: "Oakwall", Npcs: []string{}}
	call := 0
	deleted := []string{}

	provider := NewSettlementNPCProviderWithGateway(&mockNPCGateway{
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
	if len(settlement.Npcs) != 0 {
		t.Fatalf("expected settlement npc list rolled back to empty, got %v", settlement.Npcs)
	}
}

func TestSettlementNPCProvider_DeleteNPCFromSettlement_RemovesFromSettlement(t *testing.T) {
	settlement := model.Settlement{Name: "Thicket", Npcs: []string{"npc-a", "npc-b"}}
	provider := NewSettlementNPCProviderWithGateway(&mockNPCGateway{})

	err := provider.DeleteNPCFromSettlement(&settlement, "npc-a")
	if err != nil {
		t.Fatalf("DeleteNPCFromSettlement() unexpected error: %v", err)
	}
	if len(settlement.Npcs) != 1 || settlement.Npcs[0] != "npc-b" {
		t.Fatalf("expected npc-a removed from settlement list, got %v", settlement.Npcs)
	}
}
