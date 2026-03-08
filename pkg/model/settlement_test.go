package model

import "testing"

func TestSettlement_AddNpc_EmptyIDReturnsError(t *testing.T) {
	s := Settlement{}

	err := s.AddNpc("")
	if err == nil {
		t.Fatal("AddNpc() expected error for empty npc id, got nil")
	}
	if len(s.Npcs) != 0 {
		t.Fatalf("AddNpc() should not mutate NPC list on error; got %d entries", len(s.Npcs))
	}
}

func TestSettlement_AddNpc_ValidIDAppendsAndInitializesSlice(t *testing.T) {
	s := Settlement{}

	err := s.AddNpc("npc-1")
	if err != nil {
		t.Fatalf("AddNpc() unexpected error: %v", err)
	}
	if len(s.Npcs) != 1 {
		t.Fatalf("AddNpc() expected 1 npc, got %d", len(s.Npcs))
	}
	if s.Npcs[0] != "npc-1" {
		t.Fatalf("AddNpc() expected first npc id 'npc-1', got %q", s.Npcs[0])
	}
}

func TestSettlement_AddNpc_PreservesAppendOrder(t *testing.T) {
	s := Settlement{}

	for _, id := range []string{"npc-1", "npc-2", "npc-3"} {
		if err := s.AddNpc(id); err != nil {
			t.Fatalf("AddNpc() unexpected error for %q: %v", id, err)
		}
	}

	if len(s.Npcs) != 3 {
		t.Fatalf("AddNpc() expected 3 npcs, got %d", len(s.Npcs))
	}
	if s.Npcs[0] != "npc-1" || s.Npcs[1] != "npc-2" || s.Npcs[2] != "npc-3" {
		t.Fatalf("AddNpc() append order mismatch, got %v", s.Npcs)
	}
}
