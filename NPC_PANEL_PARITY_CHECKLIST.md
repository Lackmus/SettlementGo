# NPC Panel Parity Checklist

Use this checklist when changing the settlement Wails NPC panel so behavior stays aligned with NPCGen.

## Reference Files

- Settlement UI: ui/wails/dist/index.html
- Settlement behavior: ui/wails/dist/app.js
- Settlement shared helpers: ui/wails/dist/npc-shared-core.js
- Settlement styles: ui/wails/dist/style.css
- NPCGen behavior source: ../npcgen/NPCgenGo/ui/wails/dist/npc-ui-core.js
- NPCGen markup source: ../npcgen/NPCgenGo/ui/wails/dist/index.html

## Behavior Checklist

- NPC details panel and edit panel toggle exactly:
  - Details visible by default.
  - Edit opens only when an NPC is selected.
  - Close clears selection and resets form state.
- NPC roster interaction:
  - Click NPC name to load details.
  - Keyboard Enter or Space on NPC name also opens details.
  - Delete action remains available from roster card.
- Field dependencies:
  - Type change refreshes subtype list and clears stats/items placeholder.
  - Faction change refreshes species list and clears name.
  - Subtype controls reroll stats/items enablement.
  - Species controls reroll name enablement.
- Reroll behavior:
  - Reroll Stats/Items requires subtype.
  - Reroll Name requires species.
  - Error and guard messages use NPCGen-style dialog alerts.
- Save validation:
  - Required fields: name, type, subtype, species, faction, trait.
  - Missing fields show one dialog listing all missing field names.
  - Missing id shows the Generate an NPC first dialog message.
- Placeholder consistency:
  - Empty details/form display uses em dash: —

## Microcopy Lock

Keep these strings unchanged unless NPCGen changes first:

- Select an NPC first.
- Select a subtype first.
- Select a species first.
- Please fill all fields before saving. Missing: ...
- No ID present. Generate an NPC first.
- Failed to generate subtype fields.
- Failed to generate species name.
- Failed to reroll subtype fields.
- Failed to reroll name.

## Verification Steps

1. Run diagnostics for ui/wails/dist/index.html, ui/wails/dist/app.js, ui/wails/dist/style.css, ui/wails/dist/npc-shared-core.js.
2. Run: go test ./cmd/settlementgen-wails
3. Run from cmd/settlementgen-wails: wails build -clean
4. Manual smoke test:
   - Open a settlement.
   - Open NPC details from roster name.
   - Edit, reroll fields, save, cancel, close.
   - Delete one NPC and verify roster and details update.
