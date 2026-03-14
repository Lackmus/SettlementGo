package app

import (
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo"
	"github.com/lackmus/settlementgengo/internal/app/controllers"
	"github.com/lackmus/settlementgengo/internal/platform/loaders"
	"github.com/lackmus/settlementgengo/pkg/model"
	"github.com/lackmus/settlementgengo/pkg/service"
)

const (
	settlementDir  = "settlement_database"
	settlementData = "settlement_data"
	defaultDataDir = "data"
)

type SettlementGenApp struct {
	NpcGenerator               npcgengo.NPCGen
	SettlementController       *controllers.SettlementListController
	SettlementCreationSupplier *service.SettlementCreationSupplier
	SettlementService          *service.SettlementService
}

func NewSettlementGenApp() *SettlementGenApp {
	return NewSettlementGenAppWithDataDir("")
}

func NewSettlementGenAppWithDataDir(dir string) *SettlementGenApp {
	baseDir := resolveDataDir(dir)

	if err := os.MkdirAll(filepath.Join(baseDir, settlementDir), 0o755); err != nil {
		panic(err)
	}

	npcGenerator, err := npcgengo.NewNPCGenWithDataDir(baseDir)
	if err != nil {
		panic(err)
	}
	settlementService, err := service.NewSettlementService(loaders.NewJSONSettlementStorage(filepath.Join(baseDir, settlementDir)))
	if err != nil {
		panic(err)
	}
	loaders := loaders.NewJSONSettlementConfigLoader(filepath.Join(baseDir, settlementData))
	factions := npcGenerator.GetFactions()
	settlementCreationSupplier := service.NewSettlementCreationSupplier(loaders, factions)
	settlementController := controllers.NewSettlementListController(*settlementService, *settlementCreationSupplier, *npcGenerator)

	app := &SettlementGenApp{
		NpcGenerator:               *npcGenerator,
		SettlementService:          settlementService,
		SettlementCreationSupplier: settlementCreationSupplier,
		SettlementController:       settlementController,
	}
	return app
}

func resolveDataDir(dataDir string) string {
	if base := normalizeDataDir(dataDir); base != "" {
		return base
	}

	if cwd, err := os.Getwd(); err == nil {
		if base := findDataDirUp(cwd); base != "" {
			return base
		}
	}

	if executablePath, err := os.Executable(); err == nil {
		if base := findDataDirUp(filepath.Dir(executablePath)); base != "" {
			return base
		}
	}

	return defaultDataDir
}

func normalizeDataDir(base string) string {
	if base == "" {
		return ""
	}

	if hasAppData(base) {
		return base
	}

	candidate := filepath.Join(base, defaultDataDir)
	if hasAppData(candidate) {
		return candidate
	}

	return ""
}

func findDataDirUp(start string) string {
	current := start

	for {
		candidate := filepath.Join(current, defaultDataDir)
		if hasAppData(candidate) {
			return candidate
		}
		if hasAppData(current) {
			return current
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return ""
}

func hasAppData(base string) bool {
	creationData, err := os.Stat(filepath.Join(base, "creation_data", "factiondata"))
	if err != nil || !creationData.IsDir() {
		return false
	}

	settlementDataInfo, err := os.Stat(filepath.Join(base, settlementData))
	if err != nil {
		return false
	}

	return settlementDataInfo.IsDir()
}

// CreateRandomSettlementWithNPCs creates and saves a random settlement,
// then generates and attaches npcCount random NPC IDs to it.
func (a *SettlementGenApp) CreateRandomSettlementWithNPCs(npcCount int) (model.Settlement, error) {
	return a.SettlementController.CreateRandomSettlementWithNPCs(npcCount)
}

// AddRandomNPCsToSettlement generates npcCount random NPCs and appends their IDs
// to the named settlement. The updated settlement is persisted via controller update.
func (a *SettlementGenApp) AddRandomNPCsToSettlement(name string, npcCount int) (model.Settlement, error) {
	return a.SettlementController.AddRandomNPCsToSettlement(name, npcCount)
}
