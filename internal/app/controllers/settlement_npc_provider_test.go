package controllers

import (
	"errors"
	"strings"
	"testing"

	"github.com/lackmus/settlementgengo/pkg/model"
)

func TestSettlementNPCProvider_GenerateRandomNPCsForSettlement_ValidatesInit(t *testing.T) {
	provider := &SettlementNPCProvider{
		createRandomNPCFn: func() (string, error) {
			return "npc-1", nil
		},
		deleteNPCFn: func(id string) error { return nil },
	}

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

	provider := &SettlementNPCProvider{
		createRandomNPCFn: func() (string, error) {
			call++
			if call == 1 {
				return "npc-1", nil
			}
			return "", errors.New("generator failed")
		},
		deleteNPCFn: func(id string) error {
			deleted = append(deleted, id)
			return nil
		},
	}

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
	deleted := []string{}

	provider := &SettlementNPCProvider{
		deleteNPCFn: func(id string) error {
			deleted = append(deleted, id)
			return nil
		},
	}

	err := provider.DeleteNPCFromSettlement(&settlement, "npc-a")
	if err != nil {
		t.Fatalf("DeleteNPCFromSettlement() unexpected error: %v", err)
	}
	if len(deleted) != 1 || deleted[0] != "npc-a" {
		t.Fatalf("expected one delete call for npc-a, got %v", deleted)
	}
	if len(settlement.Npcs) != 1 || settlement.Npcs[0] != "npc-b" {
		t.Fatalf("expected npc-a removed from settlement list, got %v", settlement.Npcs)
	}
}
