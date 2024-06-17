package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	wdDir := filepath.Dir(os.Args[0])
	logFile, err := os.OpenFile(wdDir+"\\route-panel.log", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("open log file failed: ", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "RoutePanel",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 152, G: 152, B: 152, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
