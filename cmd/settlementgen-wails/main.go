package main

import (
	"flag"
	"log"

	settlementapp "github.com/lackmus/settlementgengo/internal/app"
	uiwails "github.com/lackmus/settlementgengo/ui/wails"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	dataDir := flag.String("data-dir", "", "path to data directory")
	flag.Parse()

	app := settlementapp.NewSettlementGenAppWithDataDir(*dataDir)
	api := NewWailsAPI(app)

	err := wails.Run(&options.App{
		Title:     "SettlementGen",
		Width:     1440,
		Height:    920,
		MinWidth:  1100,
		MinHeight: 720,
		AssetServer: &assetserver.Options{
			Assets: uiwails.Assets(),
		},
		OnStartup: api.startup,
		Bind: []interface{}{
			api,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
