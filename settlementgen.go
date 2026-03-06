package settlementgo

import "github.com/lackmus/settlementgengo/internal/app"

type SettlementGenApp = app.SettlementGenApp

func NewSettlementGenApp() *SettlementGenApp {
	return app.NewSettlementGenApp()
}

func NewSettlementGenAppWithDataDir(dir string) *SettlementGenApp {
	return app.NewSettlementGenAppWithDataDir(dir)
}
