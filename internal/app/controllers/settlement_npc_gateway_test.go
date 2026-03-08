package controllers

import (
	"strings"
	"testing"

	"github.com/lackmus/npcgengo"
)

func TestSettlementNPCGateway_DeleteNPC_RejectsEmptyID(t *testing.T) {
	gateway := newSettlementNPCGateway(npcgengo.NPCGen{})

	err := gateway.DeleteNPC("")
	if err == nil {
		t.Fatal("DeleteNPC() expected error for empty npc id, got nil")
	}
	if !strings.Contains(err.Error(), "npc id is empty") {
		t.Fatalf("expected empty id error, got: %v", err)
	}
}
