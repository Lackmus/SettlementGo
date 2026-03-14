package main

import (
	"context"
	"fmt"
	"strings"

	npcmapper "github.com/lackmus/npcgengo/pkg/mapper"
	npcmodel "github.com/lackmus/npcgengo/pkg/model"
	settlementapp "github.com/lackmus/settlementgengo/internal/app"
	"github.com/lackmus/settlementgengo/internal/app/controllers"
	appmapper "github.com/lackmus/settlementgengo/internal/app/mapper"
	settlementservice "github.com/lackmus/settlementgengo/pkg/service"
)

type SubtypeRoll struct {
	Stats string `json:"stats"`
	Items string `json:"items"`
}

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

func (a *WailsAPI) GetNPC(id string) (npcmapper.NPCInput, error) {
	controller := a.app.NpcGenerator.NPCListController
	if controller == nil {
		return npcmapper.NPCInput{}, fmt.Errorf("npc controller is not configured")
	}
	npc, err := controller.GetNPCByID(strings.TrimSpace(id))
	if err != nil {
		return npcmapper.NPCInput{}, err
	}
	return npcmapper.ToNPCInput(npc), nil
}

func (a *WailsAPI) SaveNPC(input npcmapper.NPCInput) (npcmapper.NPCInput, error) {
	controller := a.app.NpcGenerator.NPCListController
	if controller == nil {
		return npcmapper.NPCInput{}, fmt.Errorf("npc controller is not configured")
	}

	trimmedID := strings.TrimSpace(input.ID)
	if trimmedID == "" {
		return npcmapper.NPCInput{}, fmt.Errorf("cannot save without an id")
	}

	original, err := controller.GetNPCByID(trimmedID)
	if err != nil {
		return npcmapper.NPCInput{}, err
	}

	npc, err := npcmapper.ToModelNPCWithOriginal(input, controller.GetNPCBuilder(), &original)
	if err != nil {
		return npcmapper.NPCInput{}, err
	}

	if err := controller.ValidateNPC(npc); err != nil {
		return npcmapper.NPCInput{}, err
	}

	controller.UpdateNPC(npc)
	return npcmapper.ToNPCInput(npc), nil
}

func (a *WailsAPI) RollSubtypeFields(subtype string) (SubtypeRoll, error) {
	controller := a.app.NpcGenerator.NPCListController
	if controller == nil {
		return SubtypeRoll{}, fmt.Errorf("npc controller is not configured")
	}

	stats, items, err := controller.GetSubtypeFields(strings.TrimSpace(subtype))
	if err != nil {
		return SubtypeRoll{}, err
	}

	return SubtypeRoll{Stats: stats, Items: items}, nil
}

func (a *WailsAPI) RollSpeciesName(species string) (string, error) {
	controller := a.app.NpcGenerator.NPCListController
	if controller == nil {
		return "", fmt.Errorf("npc controller is not configured")
	}

	return controller.GetSpeciesName(strings.TrimSpace(species))
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

func (a *WailsAPI) UpdateSettlement(input appmapper.SettlementUpdateInput) (appmapper.SettlementView, error) {
	originalName := strings.TrimSpace(input.OriginalName)
	if originalName == "" {
		return appmapper.SettlementView{}, fmt.Errorf("original settlement name is required")
	}

	current, err := a.app.SettlementController.GetSettlement(originalName)
	if err != nil {
		return appmapper.SettlementView{}, err
	}

	updated := current
	if name := strings.TrimSpace(input.Name); name != "" {
		updated.Name = name
	}
	if faction := strings.TrimSpace(input.Faction); faction != "" {
		updated.Faction = faction
	}
	updated.Population = input.Population
	updated.Notes = strings.TrimSpace(input.Notes)

	if updated.Name == originalName {
		if err := a.app.SettlementController.UpdateSettlement(updated); err != nil {
			return appmapper.SettlementView{}, err
		}
		return a.GetSettlement(updated.Name)
	}

	if a.app.SettlementController.SettlementExists(updated.Name) {
		return appmapper.SettlementView{}, fmt.Errorf("settlement with name %q already exists", updated.Name)
	}

	if err := a.app.SettlementController.RemoveSettlement(originalName); err != nil {
		return appmapper.SettlementView{}, err
	}

	if _, err := a.app.SettlementController.AddSettlement(updated); err != nil {
		if _, rollbackErr := a.app.SettlementController.AddSettlement(current); rollbackErr != nil {
			return appmapper.SettlementView{}, fmt.Errorf("update failed: %v; rollback failed: %v", err, rollbackErr)
		}
		return appmapper.SettlementView{}, err
	}

	return a.GetSettlement(updated.Name)
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
	controller := a.app.NpcGenerator.NPCListController

	if controller == nil {
		return appmapper.ToSettlementView(settlement, nil)
	}

	return appmapper.ToSettlementView(settlement, func(id string) (npcmapper.NPCInput, error) {
		npc, err := controller.GetNPCByID(id)
		if err != nil {
			return npcmapper.NPCInput{}, err
		}
		return npcmapper.ToNPCInput(npc), nil
	})
}

func (a *WailsAPI) updateNPCInSettlements(updated npcmodel.NPC) error {
	settlements, err := a.app.SettlementController.GetAllSettlements()
	if err != nil {
		return err
	}

	for _, settlement := range settlements {
		for _, npcID := range settlement.NPCs {
			if npcID == updated.ID {
				return nil
			}
		}
	}

	return nil
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
