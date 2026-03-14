package main

import (
	"context"
	"fmt"
	"strings"

	npcmapper "github.com/lackmus/npcgengo/pkg/mapper"
	settlementapp "github.com/lackmus/settlementgengo/internal/app"
	"github.com/lackmus/settlementgengo/internal/app/controllers"
	appmapper "github.com/lackmus/settlementgengo/internal/app/mapper"
	settlementservice "github.com/lackmus/settlementgengo/pkg/service"
)

type WailsAPI struct {
	ctx context.Context
	app *settlementapp.SettlementGenApp
}

func NewWailsAPI(app *settlementapp.SettlementGenApp) *WailsAPI {
	return &WailsAPI{app: app}
}

func (a *WailsAPI) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *WailsAPI) ListSettlements() ([]appmapper.SettlementView, error) {
	settlements, err := a.app.SettlementController.GetAllSettlements()
	if err != nil {
		return nil, err
	}

	inputs := appmapper.ToSettlementInputs(settlements)
	views := make([]appmapper.SettlementView, 0, len(inputs))
	for _, input := range inputs {
		views = append(views, a.toSettlementView(input))
	}

	return views, nil
}

func (a *WailsAPI) GetSettlement(name string) (appmapper.SettlementView, error) {
	settlement, err := a.app.SettlementController.GetSettlement(strings.TrimSpace(name))
	if err != nil {
		return appmapper.SettlementView{}, err
	}

	return a.toSettlementView(appmapper.ToSettlementInput(settlement)), nil
}

func (a *WailsAPI) GetCreationOptions() controllers.CreationOptions {
	options, err := a.app.SettlementController.GetCreationOptions()
	if err != nil {
		return controllers.CreationOptions{}
	}
	return options
}

func (a *WailsAPI) CreateSettlement(input appmapper.SettlementCreateInput) (appmapper.SettlementView, error) {
	name := strings.TrimSpace(input.Name)
	faction := strings.TrimSpace(input.Faction)

	baseSettlement := settlementservice.CreateSettlement(name, faction)
	settlementInput := appmapper.ToSettlementInput(baseSettlement)
	settlementInput.XCoord = input.XCoord
	settlementInput.YCoord = input.YCoord
	if input.Population > 0 {
		settlementInput.Population = input.Population
	}
	if notes := strings.TrimSpace(input.Notes); notes != "" {
		settlementInput.Notes = notes
	}

	settlement, err := appmapper.ToSettlementModelValidated(settlementInput)
	if err != nil {
		return appmapper.SettlementView{}, err
	}

	created, err := a.app.SettlementController.AddSettlement(settlement)
	if err != nil {
		return appmapper.SettlementView{}, err
	}

	if input.InitialRandomNPCCount > 0 {
		created, err = a.app.SettlementController.AddRandomNPCsToSettlement(created.Name, input.InitialRandomNPCCount)
		if err != nil {
			return appmapper.SettlementView{}, err
		}
	}

	return a.toSettlementView(appmapper.ToSettlementInput(created)), nil
}

func (a *WailsAPI) CreateRandomSettlement() (appmapper.SettlementView, error) {
	settlement, err := a.app.SettlementController.CreateRandomSettlement()
	if err != nil {
		return appmapper.SettlementView{}, err
	}

	return a.toSettlementView(appmapper.ToSettlementInput(settlement)), nil
}

func (a *WailsAPI) CreateRandomSettlementWithNPCs(npcCount int) (appmapper.SettlementView, error) {
	settlement, err := a.app.CreateRandomSettlementWithNPCs(npcCount)
	if err != nil {
		return appmapper.SettlementView{}, err
	}

	return a.toSettlementView(appmapper.ToSettlementInput(settlement)), nil
}

func (a *WailsAPI) AddRandomNPCToSettlement(name string) (appmapper.SettlementView, error) {
	settlement, err := a.app.SettlementController.AddRandomNPCToSettlement(strings.TrimSpace(name))
	if err != nil {
		return appmapper.SettlementView{}, err
	}

	return a.toSettlementView(appmapper.ToSettlementInput(settlement)), nil
}

func (a *WailsAPI) AddRandomNPCsToSettlement(name string, npcCount int) (appmapper.SettlementView, error) {
	settlement, err := a.app.SettlementController.AddRandomNPCsToSettlement(strings.TrimSpace(name), npcCount)
	if err != nil {
		return appmapper.SettlementView{}, err
	}

	return a.toSettlementView(appmapper.ToSettlementInput(settlement)), nil
}

func (a *WailsAPI) AddNPCToSettlement(name string, npcType string, faction string) (appmapper.SettlementView, error) {
	settlement, err := a.app.SettlementController.AddNPCToSettlement(strings.TrimSpace(name), strings.TrimSpace(npcType), strings.TrimSpace(faction))
	if err != nil {
		return appmapper.SettlementView{}, err
	}

	return a.toSettlementView(appmapper.ToSettlementInput(settlement)), nil
}

func (a *WailsAPI) DeleteNPCFromSettlement(name string, npcID string) (appmapper.SettlementView, error) {
	trimmedName := strings.TrimSpace(name)
	trimmedNPCID := strings.TrimSpace(npcID)
	if err := a.deleteNPCRecords([]string{trimmedNPCID}); err != nil {
		return appmapper.SettlementView{}, err
	}
	if err := a.app.SettlementController.DeleteNPCFromSettlement(trimmedName, trimmedNPCID); err != nil {
		return appmapper.SettlementView{}, err
	}

	return a.GetSettlement(trimmedName)
}

func (a *WailsAPI) DeleteAllNPCsFromSettlement(name string) (appmapper.SettlementView, error) {
	trimmedName := strings.TrimSpace(name)
	settlement, err := a.app.SettlementController.GetSettlement(trimmedName)
	if err != nil {
		return appmapper.SettlementView{}, err
	}
	if err := a.deleteNPCRecords(settlement.NPCs); err != nil {
		return appmapper.SettlementView{}, err
	}
	if err := a.app.SettlementController.DeleteAllNPCsFromSettlement(trimmedName); err != nil {
		return appmapper.SettlementView{}, err
	}

	return a.GetSettlement(trimmedName)
}

func (a *WailsAPI) DeleteSettlement(name string) error {
	trimmedName := strings.TrimSpace(name)
	settlement, err := a.app.SettlementController.GetSettlement(trimmedName)
	if err != nil {
		return err
	}
	if err := a.deleteNPCRecords(settlement.NPCs); err != nil {
		return err
	}
	return a.app.SettlementController.RemoveSettlement(trimmedName)
}

func (a *WailsAPI) DeleteAllSettlements() error {
	settlements, err := a.app.SettlementController.GetAllSettlements()
	if err != nil {
		return err
	}

	npcIDs := make([]string, 0)
	for _, settlement := range settlements {
		npcIDs = append(npcIDs, settlement.NPCs...)
	}
	if err := a.deleteNPCRecords(npcIDs); err != nil {
		return err
	}

	return a.app.SettlementController.RemoveAllSettlements()
}

func (a *WailsAPI) toSettlementView(settlement appmapper.SettlementInputMapper) appmapper.SettlementView {
	npcs := make([]npcmapper.NPCInput, 0, len(settlement.NPCIDs))
	controller := a.app.NpcGenerator.NPCListController

	for _, npcID := range settlement.NPCIDs {
		if controller == nil {
			npcs = append(npcs, npcmapper.NPCInput{ID: npcID, Name: "NPC controller unavailable"})
			continue
		}

		npc, err := controller.GetNPCByID(npcID)
		if err != nil {
			npcs = append(npcs, npcmapper.NPCInput{
				ID:    npcID,
				Name:  "Missing NPC",
				Notes: fmt.Sprintf("Failed to load NPC: %v", err),
			})
			continue
		}

		npcs = append(npcs, npcmapper.ToNPCInput(npc))
	}

	return appmapper.SettlementView{
		Name:       settlement.Name,
		Faction:    settlement.Faction,
		XCoord:     settlement.XCoord,
		YCoord:     settlement.YCoord,
		Population: settlement.Population,
		Notes:      settlement.Notes,
		NPCs:       npcs,
	}
}

func (a *WailsAPI) deleteNPCRecords(ids []string) error {
	controller := a.app.NpcGenerator.NPCListController
	if controller == nil {
		return fmt.Errorf("npc controller is not configured")
	}

	for _, id := range ids {
		trimmedID := strings.TrimSpace(id)
		if trimmedID == "" {
			continue
		}
		controller.DeleteNPC(trimmedID)
	}

	return nil
}
