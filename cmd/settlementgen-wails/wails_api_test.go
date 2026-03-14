package main

import (
	"os"
	"path/filepath"
	"testing"

	settlementapp "github.com/lackmus/settlementgengo/internal/app"
	appmapper "github.com/lackmus/settlementgengo/internal/app/mapper"
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

		data, err := os.ReadFile(srcPath)
		if err != nil {
			t.Fatalf("failed to read %s: %v", srcPath, err)
		}
		if err := os.WriteFile(dstPath, data, 0o644); err != nil {
			t.Fatalf("failed to write %s: %v", dstPath, err)
		}
	}
}

func newWailsAPIForTests(t *testing.T) *WailsAPI {
	t.Helper()

	baseDir := t.TempDir()
	copyDir(t, filepath.Clean("../../data/creation_data"), filepath.Join(baseDir, "creation_data"))
	copyDir(t, filepath.Clean("../../data/settlement_data"), filepath.Join(baseDir, "settlement_data"))

	if err := os.MkdirAll(filepath.Join(baseDir, "npc_database"), 0o755); err != nil {
		t.Fatalf("failed to create npc database dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(baseDir, "settlement_database"), 0o755); err != nil {
		t.Fatalf("failed to create settlement database dir: %v", err)
	}

	app := settlementapp.NewSettlementGenAppWithDataDir(baseDir)
	t.Cleanup(func() {
		_ = app.SettlementController.RemoveAllSettlements()
		app.NpcGenerator.NPCListController.DeleteAllNPCs()
	})

	return NewWailsAPI(app)
}

func TestWailsAPI_GetCreationOptions_ReturnsNPCOptions(t *testing.T) {
	api := newWailsAPIForTests(t)

	options := api.GetCreationOptions()
	if len(options.Factions) == 0 {
		t.Fatal("expected factions to be available")
	}
	if len(options.NpcTypes) == 0 {
		t.Fatal("expected npc types to be available")
	}
	if len(options.NpcSubtypeForTypeMap) == 0 {
		t.Fatal("expected npc subtype map to be available")
	}
}

func TestWailsAPI_CreateRandomSettlementWithNPCs_ResolvesNPCs(t *testing.T) {
	api := newWailsAPIForTests(t)

	settlement, err := api.CreateRandomSettlementWithNPCs(2)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if settlement.Name == "" {
		t.Fatal("expected settlement name to be populated")
	}
	if len(settlement.NPCs) != 2 {
		t.Fatalf("expected 2 NPCs, got %d", len(settlement.NPCs))
	}
	for _, npc := range settlement.NPCs {
		if npc.ID == "" {
			t.Fatal("expected generated npc id to be populated")
		}
	}

	settlements, err := api.ListSettlements()
	if err != nil {
		t.Fatalf("expected list call to succeed, got: %v", err)
	}
	if len(settlements) != 1 {
		t.Fatalf("expected 1 stored settlement, got %d", len(settlements))
	}
}

func TestWailsAPI_AddAndDeleteNPCFromSettlement(t *testing.T) {
	api := newWailsAPIForTests(t)

	created, err := api.CreateSettlement(appmapper.SettlementCreateInput{
		Name:       "Iron Hollow",
		Faction:    api.GetCreationOptions().Factions[0],
		Population: 240,
		Notes:      "Border outpost",
	})
	if err != nil {
		t.Fatalf("expected settlement creation to succeed, got: %v", err)
	}

	updated, err := api.AddRandomNPCToSettlement(created.Name)
	if err != nil {
		t.Fatalf("expected random npc addition to succeed, got: %v", err)
	}
	if len(updated.NPCs) != 1 {
		t.Fatalf("expected 1 NPC after add, got %d", len(updated.NPCs))
	}

	updated, err = api.DeleteNPCFromSettlement(created.Name, updated.NPCs[0].ID)
	if err != nil {
		t.Fatalf("expected npc deletion to succeed, got: %v", err)
	}
	if len(updated.NPCs) != 0 {
		t.Fatalf("expected no NPCs after delete, got %d", len(updated.NPCs))
	}
	if updated.Notes != "Border outpost" {
		t.Fatalf("expected notes to remain unchanged, got %q", updated.Notes)
	}
}

func TestWailsAPI_UpdateSettlement_WithoutCoords(t *testing.T) {
	api := newWailsAPIForTests(t)

	options := api.GetCreationOptions()
	if len(options.Factions) < 2 {
		t.Fatal("expected at least two factions for update test")
	}

	created, err := api.CreateSettlement(appmapper.SettlementCreateInput{
		Name:       "Mossfield",
		Faction:    options.Factions[0],
		XCoord:     111,
		YCoord:     222,
		Population: 320,
		Notes:      "Old roads and farms",
	})
	if err != nil {
		t.Fatalf("expected settlement creation to succeed, got: %v", err)
	}

	withNPC, err := api.AddRandomNPCToSettlement(created.Name)
	if err != nil {
		t.Fatalf("expected random npc addition to succeed, got: %v", err)
	}
	if len(withNPC.NPCs) != 1 {
		t.Fatalf("expected 1 NPC after add, got %d", len(withNPC.NPCs))
	}

	updated, err := api.UpdateSettlement(appmapper.SettlementUpdateInput{
		OriginalName: created.Name,
		Name:         "Mossfield Prime",
		Faction:      options.Factions[1],
		Population:   415,
		Notes:        "Expanded trade quarter",
	})
	if err != nil {
		t.Fatalf("expected settlement update to succeed, got: %v", err)
	}

	if updated.Name != "Mossfield Prime" {
		t.Fatalf("expected renamed settlement, got %q", updated.Name)
	}
	if updated.Faction != options.Factions[1] {
		t.Fatalf("expected faction %q, got %q", options.Factions[1], updated.Faction)
	}
	if updated.Population != 415 {
		t.Fatalf("expected population 415, got %d", updated.Population)
	}
	if updated.Notes != "Expanded trade quarter" {
		t.Fatalf("expected updated notes, got %q", updated.Notes)
	}
	if updated.XCoord != 111 || updated.YCoord != 222 {
		t.Fatalf("expected coordinates unchanged, got (%d, %d)", updated.XCoord, updated.YCoord)
	}
	if len(updated.NPCs) != 1 {
		t.Fatalf("expected npc links preserved, got %d npcs", len(updated.NPCs))
	}

	if _, err := api.GetSettlement(created.Name); err == nil {
		t.Fatal("expected old settlement name to be unavailable after rename")
	}
}
