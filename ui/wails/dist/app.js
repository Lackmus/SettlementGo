import { textOrDash, isPresent, validatePayload } from "./npc-shared-core.js";

function getBackend() {
  const api = window?.go?.main?.WailsAPI;
  if (!api) {
    throw new Error("Wails bindings are unavailable. Start the app with Wails to use the desktop UI.");
  }
  return api;
}

const state = {
  settlements: [],
  selectedName: "",
  editingSettlement: false,
  selectedNpcId: "",
  selectedNpcSnapshot: null,
  creationOptions: {
    factions: [],
    npcTypes: [],
    npcSubtypeForTypeMap: {},
    npcSpeciesForFactionMap: {},
    species: [],
    traits: [],
  },
};

const elements = {
  customForm: document.querySelector("#customForm"),
  randomForm: document.querySelector("#randomForm"),
  refreshButton: document.querySelector("#refreshButton"),
  clearAllButton: document.querySelector("#clearAllButton"),
  settlementList: document.querySelector("#settlementList"),
  npcList: document.querySelector("#npcList"),
  message: document.querySelector("#message"),
  emptyState: document.querySelector("#emptyState"),
  detailContent: document.querySelector("#detailContent"),
  detailTitle: document.querySelector("#detailTitle"),
  detailFaction: document.querySelector("#detailFaction"),
  detailPopulation: document.querySelector("#detailPopulation"),
  detailNpcCount: document.querySelector("#detailNpcCount"),
  detailNotes: document.querySelector("#detailNotes"),
  editSettlementButton: document.querySelector("#editSettlementButton"),
  settlementEditForm: document.querySelector("#settlementEditForm"),
  editSettlementName: document.querySelector("#editSettlementName"),
  editSettlementFaction: document.querySelector("#editSettlementFaction"),
  editSettlementPopulation: document.querySelector("#editSettlementPopulation"),
  editSettlementNotes: document.querySelector("#editSettlementNotes"),
  cancelSettlementEditButton: document.querySelector("#cancelSettlementEditButton"),
  deleteSettlementButton: document.querySelector("#deleteSettlementButton"),
  purgeSettlementNpcsButton: document.querySelector("#purgeSettlementNpcsButton"),
  addRandomNpcButton: document.querySelector("#addRandomNpcButton"),
  specificNpcForm: document.querySelector("#specificNpcForm"),
  settlementFaction: document.querySelector("#settlementFaction"),
  npcFactionSelect: document.querySelector("#npcFactionSelect"),
  npcTypeSelect: document.querySelector("#npcTypeSelect"),
  statSettlements: document.querySelector("#statSettlements"),
  statNpcs: document.querySelector("#statNpcs"),
  statFaction: document.querySelector("#statFaction"),
  npcDetails: document.querySelector("#npcDetails"),
  btnCloseNpcDetail: document.querySelector("#btnCloseNpcDetail"),
  btnEditNpcDetail: document.querySelector("#btnEditNpcDetail"),
  npcForm: document.querySelector("#npcForm"),
  fId: document.querySelector("#f_id"),
  fName: document.querySelector("#f_name"),
  fType: document.querySelector("#f_type"),
  fSubtype: document.querySelector("#f_subtype"),
  fSpecies: document.querySelector("#f_species"),
  fFaction: document.querySelector("#f_faction"),
  fTraits: document.querySelector("#f_traits"),
  fStats: document.querySelector("#f_stats"),
  fItems: document.querySelector("#f_items"),
  fNotes: document.querySelector("#f_notes"),
  btnReroll: document.querySelector("#btnReroll"),
  btnRerollName: document.querySelector("#btnRerollName"),
  btnSaveNpc: document.querySelector("#btnSaveNpc"),
  btnCancelNpcEdit: document.querySelector("#btnCancelNpcEdit"),
  dName: document.querySelector("#d_name"),
  dType: document.querySelector("#d_type"),
  dSubtype: document.querySelector("#d_subtype"),
  dSpecies: document.querySelector("#d_species"),
  dFaction: document.querySelector("#d_faction"),
  dTraits: document.querySelector("#d_traits"),
  dStats: document.querySelector("#d_stats"),
  dItems: document.querySelector("#d_items"),
  dNotes: document.querySelector("#d_notes"),
};

function setSettlementEditMode(enabled) {
  state.editingSettlement = enabled;
  elements.settlementEditForm.classList.toggle("hidden", !enabled);
  elements.editSettlementButton.classList.toggle("hidden", enabled);
}

function readSettlementEditForm() {
  return {
    name: elements.editSettlementName.value.trim(),
    faction: elements.editSettlementFaction.value.trim(),
    population: Number(elements.editSettlementPopulation.value || 0),
    notes: elements.editSettlementNotes.value.trim(),
  };
}

function isSettlementEditDirty(settlement) {
  if (!state.editingSettlement || !settlement) {
    return false;
  }

  const current = readSettlementEditForm();
  return (
    current.name !== String(settlement.name || "").trim() ||
    current.faction !== String(settlement.faction || "").trim() ||
    current.population !== Number(settlement.population ?? 0) ||
    current.notes !== String(settlement.notes || "").trim()
  );
}

function confirmDiscardSettlementEdits(settlement) {
  if (!isSettlementEditDirty(settlement)) {
    return true;
  }
  return window.confirm("Discard unsaved settlement changes?");
}

function populateSettlementEditForm(settlement) {
  if (!settlement) {
    elements.editSettlementName.value = "";
    elements.editSettlementFaction.value = "";
    elements.editSettlementPopulation.value = "0";
    elements.editSettlementNotes.value = "";
    return;
  }

  elements.editSettlementName.value = settlement.name || "";
  setSelectValue(elements.editSettlementFaction, settlement.faction || "");
  elements.editSettlementPopulation.value = String(settlement.population ?? 0);
  elements.editSettlementNotes.value = settlement.notes || "";
}

function selectedNpc() {
  const settlement = selectedSettlement();
  if (!settlement) {
    return null;
  }
  return settlement.npcs.find((npc) => npc.id === state.selectedNpcId) || null;
}

function clearNpcDetailSelection() {
  state.selectedNpcId = "";
  state.selectedNpcSnapshot = null;
}

function clearNpcFormState() {
  elements.fId.value = "";
  elements.fName.value = "";
  setSelectValue(elements.fType, "");
  setSelectValue(elements.fSubtype, "");
  setSelectValue(elements.fSpecies, "");
  setSelectValue(elements.fFaction, "");
  setSelectValue(elements.fTraits, "");
  elements.fStats.textContent = "—";
  elements.fItems.textContent = "—";
  elements.fNotes.value = "";
  setFieldEnabled(elements.fSubtype, false);
  setFieldEnabled(elements.fSpecies, false);
  setButtonEnabled(elements.btnReroll, false);
  setButtonEnabled(elements.btnRerollName, false);
}

function setButtonEnabled(button, enabled) {
  if (!button) {
    return;
  }
  button.disabled = !enabled;
}

function setFieldEnabled(field, enabled) {
  if (!field) {
    return;
  }
  field.disabled = !enabled;
}

function setSelectOptions(select, values, includeEmpty = true) {
  if (!select) {
    return;
  }

  const currentValue = select.value;
  select.innerHTML = "";

  if (includeEmpty) {
    const empty = document.createElement("option");
    empty.value = "";
    empty.textContent = "";
    select.append(empty);
  }

  for (const value of values || []) {
    const option = document.createElement("option");
    option.value = value;
    option.textContent = value;
    select.append(option);
  }

  if (currentValue && Array.from(select.options).some((option) => option.value === currentValue)) {
    select.value = currentValue;
  }
}

function setSelectValue(select, value) {
  if (!select) {
    return;
  }
  const desired = value || "";
  if (desired && !Array.from(select.options).some((option) => option.value === desired)) {
    const option = document.createElement("option");
    option.value = desired;
    option.textContent = desired;
    select.append(option);
  }
  select.value = desired;
}

function showNpcDetailsPanel() {
  elements.npcDetails.classList.remove("hidden");
  elements.npcForm.classList.add("hidden");
}

function showNpcEditPanel() {
  elements.npcDetails.classList.add("hidden");
  elements.npcForm.classList.remove("hidden");
}

function updateSubtypeDropdown(selectedType, selectedSubtype = "") {
  const subtypeMap = state.creationOptions.npcSubtypeForTypeMap || {};
  const subtypes = subtypeMap[selectedType] || [];
  setSelectOptions(elements.fSubtype, subtypes, true);
  setSelectValue(elements.fSubtype, selectedSubtype);
  setFieldEnabled(elements.fSubtype, isPresent(selectedType));
  setButtonEnabled(elements.btnReroll, isPresent(elements.fSubtype.value));
}

function updateSpeciesDropdown(selectedFaction, selectedSpecies = "") {
  const speciesMap = state.creationOptions.npcSpeciesForFactionMap || {};
  const species = speciesMap[selectedFaction] || [];
  setSelectOptions(elements.fSpecies, species, true);
  setSelectValue(elements.fSpecies, selectedSpecies);
  setFieldEnabled(elements.fSpecies, isPresent(selectedFaction));
  setButtonEnabled(elements.btnRerollName, isPresent(elements.fSpecies.value));
}

function setNpcForm(npc) {
  state.selectedNpcSnapshot = npc ? { ...npc } : null;
  elements.fId.value = npc?.id || "";
  elements.fName.value = npc?.name || "";
  setSelectValue(elements.fType, npc?.type || "");
  updateSubtypeDropdown(npc?.type || "", npc?.subtype || "");
  setSelectValue(elements.fFaction, npc?.faction || "");
  updateSpeciesDropdown(npc?.faction || "", npc?.species || "");
  const trait = (npc?.trait || "").split(",")[0]?.trim() || "";
  setSelectValue(elements.fTraits, trait);
  elements.fStats.textContent = npc?.stats || "—";
  elements.fItems.textContent = npc?.items || "—";
  elements.fNotes.value = npc?.notes || "";

  setButtonEnabled(elements.btnReroll, isPresent(elements.fSubtype.value));
  setButtonEnabled(elements.btnRerollName, isPresent(elements.fSpecies.value));
}

function readNpcForm() {
  const statsValue = (elements.fStats.textContent || "").trim();
  const itemsValue = (elements.fItems.textContent || "").trim();
  return {
    id: elements.fId.value || "",
    name: elements.fName.value || "",
    type: elements.fType.value || "",
    subtype: elements.fSubtype.value || "",
    species: elements.fSpecies.value || "",
    faction: elements.fFaction.value || "",
    trait: elements.fTraits.value || "",
    stats: statsValue === "—" ? "" : statsValue,
    items: itemsValue === "—" ? "" : itemsValue,
    notes: elements.fNotes.value || "",
  };
}

async function applySubtypeRoll(subtype) {
  if (!isPresent(subtype)) {
    elements.fStats.textContent = "—";
    elements.fItems.textContent = "—";
    return;
  }
  const rolled = await getBackend().RollSubtypeFields(subtype);
  elements.fStats.textContent = rolled?.stats || rolled?.Stats || "—";
  elements.fItems.textContent = rolled?.items || rolled?.Items || "—";
}

async function applySpeciesNameRoll(species) {
  if (!isPresent(species)) {
    elements.fName.value = "";
    return;
  }
  const rolledName = await getBackend().RollSpeciesName(species);
  elements.fName.value = rolledName || "";
}

function renderNpcDetails() {
  const npc = selectedNpc();
  if (!npc) {
    showNpcDetailsPanel();
    elements.dName.textContent = "—";
    elements.dType.textContent = "—";
    elements.dSubtype.textContent = "—";
    elements.dSpecies.textContent = "—";
    elements.dFaction.textContent = "—";
    elements.dTraits.textContent = "—";
    elements.dStats.textContent = "—";
    elements.dItems.textContent = "—";
    elements.dNotes.textContent = "—";
    elements.npcForm.classList.add("hidden");
    setButtonEnabled(elements.btnEditNpcDetail, false);
    return;
  }

  showNpcDetailsPanel();
  setButtonEnabled(elements.btnEditNpcDetail, isPresent(npc.id));
  elements.dName.textContent = textOrDash(npc.name);
  elements.dType.textContent = textOrDash(npc.type);
  elements.dSubtype.textContent = textOrDash(npc.subtype);
  elements.dSpecies.textContent = textOrDash(npc.species);
  elements.dFaction.textContent = textOrDash(npc.faction);
  elements.dTraits.textContent = textOrDash(npc.trait);
  elements.dStats.textContent = textOrDash(npc.stats);
  elements.dItems.textContent = textOrDash(npc.items);
  elements.dNotes.textContent = textOrDash(npc.notes);
}

function selectedSettlement() {
  return state.settlements.find((settlement) => settlement.name === state.selectedName) || null;
}

function showMessage(text, type = "info") {
  if (!text) {
    elements.message.classList.add("hidden");
    elements.message.textContent = "";
    elements.message.classList.remove("error");
    return;
  }

  elements.message.textContent = text;
  elements.message.classList.remove("hidden");
  elements.message.classList.toggle("error", type === "error");
}

function populateSelect(select, values, placeholder) {
  const normalized = Array.from(new Set(values.filter(Boolean))).sort((left, right) => left.localeCompare(right));
  select.innerHTML = "";

  if (placeholder) {
    const option = document.createElement("option");
    option.value = "";
    option.textContent = placeholder;
    select.append(option);
  }

  normalized.forEach((value) => {
    const option = document.createElement("option");
    option.value = value;
    option.textContent = value;
    select.append(option);
  });

  if (!placeholder && normalized.length > 0) {
    select.value = normalized[0];
  }
}

function renderStats() {
  const totalNpcCount = state.settlements.reduce((count, settlement) => count + settlement.npcs.length, 0);
  const selected = selectedSettlement();

  elements.statSettlements.textContent = String(state.settlements.length);
  elements.statNpcs.textContent = String(totalNpcCount);
  elements.statFaction.textContent = selected?.faction || state.settlements[0]?.faction || "-";
}

function renderSettlementList() {
  elements.settlementList.innerHTML = "";

  if (state.settlements.length === 0) {
    const empty = document.createElement("div");
    empty.className = "empty-state";
    empty.innerHTML = "<div><h3>No settlements yet</h3><p>Generate one to populate the registry.</p></div>";
    elements.settlementList.append(empty);
    return;
  }

  state.settlements.forEach((settlement) => {
    const card = document.createElement("article");
    card.className = "settlement-card";
    if (settlement.name === state.selectedName) {
      card.classList.add("active");
    }

    const faction = settlement.faction || "Unknown faction";
    const notes = settlement.notes || "No notes recorded.";

    card.innerHTML = `
      <button type="button" class="select-card">
        <div class="card-header">
          <div>
            <h3>${settlement.name}</h3>
            <p class="eyebrow">${faction}</p>
          </div>
          <strong>${settlement.npcs.length} NPCs</strong>
        </div>
        <div class="pill-row">
          <span class="pill">Population ${settlement.population}</span>
        </div>
        <p class="meta-line">${notes}</p>
      </button>
    `;

    card.querySelector("button").addEventListener("click", async () => {
      if (!confirmDiscardSettlementEdits(selectedSettlement())) {
        return;
      }
      state.selectedName = settlement.name;
      setSettlementEditMode(false);
      clearNpcDetailSelection();
      await refreshSelection();
    });

    elements.settlementList.append(card);
  });
}

function renderDetailPane() {
  const settlement = selectedSettlement();
  const hasSelection = Boolean(settlement);

  elements.emptyState.classList.toggle("hidden", hasSelection);
  elements.detailContent.classList.toggle("hidden", !hasSelection);
  elements.deleteSettlementButton.disabled = !hasSelection;
  elements.purgeSettlementNpcsButton.disabled = !hasSelection;
  elements.addRandomNpcButton.disabled = !hasSelection;
  elements.editSettlementButton.disabled = !hasSelection;
  elements.specificNpcForm.querySelector("button").disabled = !hasSelection;

  if (!settlement) {
    setSettlementEditMode(false);
    populateSettlementEditForm(null);
    elements.detailTitle.textContent = "Select a settlement";
    elements.npcList.innerHTML = "";
    elements.npcDetails.classList.add("hidden");
    elements.npcForm.classList.add("hidden");
    return;
  }

  elements.detailTitle.textContent = settlement.name;
  elements.detailFaction.textContent = settlement.faction || "-";
  elements.detailPopulation.textContent = String(settlement.population ?? 0);
  elements.detailNpcCount.textContent = String(settlement.npcs.length);
  elements.detailNotes.textContent = settlement.notes || "No notes recorded.";
  populateSettlementEditForm(settlement);

  elements.npcList.innerHTML = "";
  if (settlement.npcs.length === 0) {
    clearNpcDetailSelection();
    renderNpcDetails();
    const empty = document.createElement("div");
    empty.className = "empty-state";
    empty.innerHTML = "<div><h3>No NPCs assigned</h3><p>Add a random or specific NPC to populate this settlement.</p></div>";
    elements.npcList.append(empty);
    return;
  }

  settlement.npcs.forEach((npc) => {
    const card = document.createElement("article");
    card.className = "npc-card";
    if (state.selectedNpcId === npc.id) {
      card.classList.add("active");
    }
    card.innerHTML = `
      <div class="npc-header">
        <div>
          <h3 class="npc-name" role="button" tabindex="0">${npc.name || npc.id}</h3>
          <span class="meta">${npc.type || "Unknown type"} • ${npc.subtype || "Unknown subtype"}</span>
        </div>
        <span class="pill">${npc.species || "Unknown species"}</span>
      </div>
      <p class="meta-line"><strong>Faction:</strong> ${npc.faction || "-"}</p>
      <div class="button-row">
        <button type="button" class="ghost danger">Delete</button>
      </div>
    `;

    const nameButton = card.querySelector(".npc-name");
    const deleteButton = card.querySelector("button");

    const openDetails = async () => {
      const loaded = await getBackend().GetNPC(npc.id);
      state.selectedNpcId = loaded.id || npc.id;
      setNpcForm(loaded);
      renderDetailPane();
    };

    nameButton.addEventListener("click", openDetails);
    nameButton.addEventListener("keydown", async (event) => {
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        await openDetails();
      }
    });

    deleteButton.addEventListener("click", async () => {
      if (!window.confirm(`Delete NPC ${npc.name || npc.id} from ${settlement.name}?`)) {
        return;
      }

      await mutate(async () => getBackend().DeleteNPCFromSettlement(settlement.name, npc.id), `${npc.name || npc.id} deleted.`);
    });

    elements.npcList.append(card);
  });

  if (!selectedNpc()) {
    state.selectedNpcId = settlement.npcs[0]?.id || "";
  }
  if (selectedNpc()) {
    setNpcForm(selectedNpc());
  }
  renderNpcDetails();
}

async function loadCreationOptions() {
  state.creationOptions = await getBackend().GetCreationOptions();
  populateSelect(elements.settlementFaction, state.creationOptions.factions, "Choose a faction");
  populateSelect(elements.editSettlementFaction, state.creationOptions.factions, "Choose a faction");
  populateSelect(elements.npcFactionSelect, state.creationOptions.factions, "Match settlement faction");
  populateSelect(elements.npcTypeSelect, state.creationOptions.npcTypes, "Choose an NPC type");
  setSelectOptions(elements.fType, state.creationOptions.npcTypes, true);
  setSelectOptions(elements.fFaction, state.creationOptions.factions, true);
  setSelectOptions(elements.fTraits, state.creationOptions.traits, true);
  updateSubtypeDropdown("");
  updateSpeciesDropdown("");

  setFieldEnabled(elements.fSubtype, false);
  setFieldEnabled(elements.fSpecies, false);
  setButtonEnabled(elements.btnReroll, false);
  setButtonEnabled(elements.btnRerollName, false);
  setButtonEnabled(elements.btnEditNpcDetail, false);
}

async function loadSettlements() {
  state.settlements = await getBackend().ListSettlements();

  if (state.selectedName) {
    const exists = state.settlements.some((settlement) => settlement.name === state.selectedName);
    if (!exists) {
      state.selectedName = state.settlements[0]?.name || "";
    }
  } else {
    state.selectedName = state.settlements[0]?.name || "";
  }
}

async function refreshSelection() {
  if (!state.selectedName) {
    renderStats();
    renderSettlementList();
    renderDetailPane();
    return;
  }

  const updated = await getBackend().GetSettlement(state.selectedName);
  state.settlements = state.settlements.map((settlement) => settlement.name === updated.name ? updated : settlement);
  renderStats();
  renderSettlementList();
  renderDetailPane();
}

async function syncUI() {
  await loadSettlements();
  renderStats();
  renderSettlementList();
  renderDetailPane();
}

async function mutate(action, successMessage) {
  showMessage("");

  try {
    const result = await action();
    if (result?.name) {
      const index = state.settlements.findIndex((settlement) => settlement.name === result.name);
      if (index >= 0) {
        state.settlements[index] = result;
      } else {
        state.settlements.unshift(result);
      }
      state.selectedName = result.name;
    }

    await syncUI();
    showMessage(successMessage);
  } catch (error) {
    console.error(error);
    showMessage(error?.message || String(error), "error");
  }
}

function bindEvents() {
  elements.randomForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const npcCount = Number(document.querySelector("#randomNpcCount").value || 0);
    await mutate(async () => getBackend().CreateRandomSettlementWithNPCs(npcCount), "Random settlement created.");
  });

  elements.customForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const payload = {
      name: document.querySelector("#settlementName").value.trim(),
      faction: elements.settlementFaction.value.trim(),
      population: Number(document.querySelector("#settlementPopulation").value || 0),
      initialRandomNpcCount: Number(document.querySelector("#settlementNpcCount").value || 0),
      xCoord: 0,
      yCoord: 0,
      notes: document.querySelector("#settlementNotes").value.trim(),
    };

    await mutate(async () => getBackend().CreateSettlement(payload), `Settlement ${payload.name} created.`);
    setSettlementEditMode(false);
    elements.customForm.reset();
    if (state.creationOptions.factions[0]) {
      elements.settlementFaction.value = state.creationOptions.factions[0];
      elements.npcFactionSelect.value = state.creationOptions.factions[0];
    }
  });

  elements.refreshButton.addEventListener("click", async () => {
    if (!confirmDiscardSettlementEdits(selectedSettlement())) {
      return;
    }
    showMessage("");
    setSettlementEditMode(false);
    await syncUI();
  });

  elements.editSettlementButton.addEventListener("click", () => {
    const settlement = selectedSettlement();
    if (!settlement) {
      return;
    }
    populateSettlementEditForm(settlement);
    setSettlementEditMode(true);
  });

  elements.cancelSettlementEditButton.addEventListener("click", () => {
    if (!confirmDiscardSettlementEdits(selectedSettlement())) {
      return;
    }
    populateSettlementEditForm(selectedSettlement());
    setSettlementEditMode(false);
  });

  elements.settlementEditForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const settlement = selectedSettlement();
    if (!settlement) {
      return;
    }

    const payload = {
      originalName: settlement.name,
      ...readSettlementEditForm(),
    };

    if (!isPresent(payload.name)) {
      window.alert("Settlement name cannot be empty.");
      return;
    }
    if (!isPresent(payload.faction)) {
      window.alert("Faction cannot be empty.");
      return;
    }
    if (payload.population < 0) {
      window.alert("Population cannot be negative.");
      return;
    }

    await mutate(async () => getBackend().UpdateSettlement(payload), `Settlement ${payload.name} updated.`);
    setSettlementEditMode(false);
  });

  elements.clearAllButton.addEventListener("click", async () => {
    if (!window.confirm("Delete every stored settlement?")) {
      return;
    }

    await mutate(async () => {
      await getBackend().DeleteAllSettlements();
      state.selectedName = "";
      setSettlementEditMode(false);
      clearNpcDetailSelection();
      return null;
    }, "All settlements removed.");
  });

  elements.deleteSettlementButton.addEventListener("click", async () => {
    const settlement = selectedSettlement();
    if (!settlement) {
      return;
    }
    if (!window.confirm(`Delete settlement ${settlement.name}?`)) {
      return;
    }

    await mutate(async () => {
      await getBackend().DeleteSettlement(settlement.name);
      state.selectedName = "";
      setSettlementEditMode(false);
      clearNpcDetailSelection();
      return null;
    }, `${settlement.name} deleted.`);
  });

  elements.purgeSettlementNpcsButton.addEventListener("click", async () => {
    const settlement = selectedSettlement();
    if (!settlement) {
      return;
    }
    if (!window.confirm(`Delete all NPC records attached to ${settlement.name}?`)) {
      return;
    }

    await mutate(async () => getBackend().DeleteAllNPCsFromSettlement(settlement.name), `Deleted settlement NPCs from ${settlement.name}.`);
  });

  elements.addRandomNpcButton.addEventListener("click", async () => {
    const settlement = selectedSettlement();
    if (!settlement) {
      return;
    }

    await mutate(async () => getBackend().AddRandomNPCToSettlement(settlement.name), `Random NPC added to ${settlement.name}.`);
  });

  elements.specificNpcForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const settlement = selectedSettlement();
    if (!settlement) {
      return;
    }

    const npcType = elements.npcTypeSelect.value;
    const faction = elements.npcFactionSelect.value || settlement.faction;
    await mutate(async () => getBackend().AddNPCToSettlement(settlement.name, npcType, faction), `Specific NPC added to ${settlement.name}.`);
  });

  elements.btnCloseNpcDetail.addEventListener("click", () => {
    clearNpcDetailSelection();
    clearNpcFormState();
    renderNpcDetails();
  });

  elements.btnEditNpcDetail.addEventListener("click", () => {
    if (!selectedNpc()) {
      window.alert("Select an NPC first.");
      return;
    }
    showNpcEditPanel();
  });

  elements.fType.addEventListener("change", () => {
    updateSubtypeDropdown(elements.fType.value);
    elements.fStats.textContent = "—";
    elements.fItems.textContent = "—";
  });

  elements.fFaction.addEventListener("change", () => {
    updateSpeciesDropdown(elements.fFaction.value);
    elements.fName.value = "";
  });

  elements.fSubtype.addEventListener("change", async () => {
    setButtonEnabled(elements.btnReroll, isPresent(elements.fSubtype.value));
    try {
      await applySubtypeRoll(elements.fSubtype.value);
    } catch (error) {
      window.alert(error?.message || "Failed to generate subtype fields.");
    }
  });

  elements.fSpecies.addEventListener("change", async () => {
    setButtonEnabled(elements.btnRerollName, isPresent(elements.fSpecies.value));
    try {
      await applySpeciesNameRoll(elements.fSpecies.value);
    } catch (error) {
      window.alert(error?.message || "Failed to generate species name.");
    }
  });

  elements.btnReroll.addEventListener("click", async () => {
    if (!isPresent(elements.fSubtype.value)) {
      window.alert("Select a subtype first.");
      return;
    }
    try {
      await applySubtypeRoll(elements.fSubtype.value);
    } catch (error) {
      window.alert(error?.message || "Failed to reroll subtype fields.");
    }
  });

  elements.btnRerollName.addEventListener("click", async () => {
    if (!isPresent(elements.fSpecies.value)) {
      window.alert("Select a species first.");
      return;
    }
    try {
      await applySpeciesNameRoll(elements.fSpecies.value);
    } catch (error) {
      window.alert(error?.message || "Failed to reroll name.");
    }
  });

  elements.btnSaveNpc.addEventListener("click", async () => {
    const payload = readNpcForm();
    const validation = validatePayload(payload);
    if (!validation.ok) {
      window.alert(`Please fill all fields before saving. Missing: ${validation.missing.join(", ")}`);
      return;
    }
    if (!isPresent(payload.id)) {
      window.alert("No ID present. Generate an NPC first.");
      return;
    }

    await mutate(async () => {
      const saved = await getBackend().SaveNPC(payload);
      state.selectedNpcId = saved.id || payload.id;
      return await getBackend().GetSettlement(state.selectedName);
    }, "NPC updated.");
    const latest = selectedNpc();
    if (latest) {
      setNpcForm(latest);
    }
    renderNpcDetails();
  });

  elements.btnCancelNpcEdit.addEventListener("click", () => {
    const fallback = state.selectedNpcSnapshot || selectedNpc();
    if (fallback) {
      setNpcForm(fallback);
      renderNpcDetails();
    } else {
      clearNpcDetailSelection();
      renderNpcDetails();
    }
  });
}

async function main() {
  try {
    await loadCreationOptions();
    bindEvents();
    await syncUI();
  } catch (error) {
    console.error(error);
    showMessage(error?.message || "Failed to start SettlementGen UI.", "error");
  }
}

main();