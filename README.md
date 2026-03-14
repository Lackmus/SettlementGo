# SettlementGen

Desktop settlement generator with NPC integration via npcgengo.

## Run And Build (Wails UI)

### Recommended (Wails CLI)

From cmd/settlementgen-wails:

```powershell
wails dev
```

```powershell
wails build -clean
```

Output executable:

- cmd/settlementgen-wails/build/bin/settlementgen-wails.exe

### VS Code Tasks

Use these tasks from the Command Palette (Run Task):

- Wails: Dev (SettlementGen)
- Wails: Build (SettlementGen)
- Go: Run Wails with dev tag
- Go: Build Wails with production tag

### Manual Go Commands (If Not Using Wails CLI)

Wails requires build tags for direct go run/go build commands.

```powershell
go run -tags dev ./cmd/settlementgen-wails
```

```powershell
go build -tags production ./cmd/settlementgen-wails
```

If you build without these tags, Wails will show the runtime error dialog about incorrect build tags.
