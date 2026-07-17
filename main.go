package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/Enigamitsuj/among-us-mod-launcher/internal/appservice"
	"github.com/Enigamitsuj/among-us-mod-launcher/internal/mods"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[mods.InstallProgress]("install:progress")
}

func main() {
	launcher := appservice.New()

	app := application.New(application.Options{
		Name:        "Among Us Mod Launcher",
		Description: "Community-made launcher for Among Us mods. Installer created by FBI OpenUp.",
		Services: []application.Service{
			application.NewService(launcher),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	const width, height = 1220, 780

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:             "Among Us Mod Launcher",
		Width:             width,
		Height:            height,
		MinWidth:          width,
		MinHeight:         height,
		MaxWidth:          width,
		MaxHeight:         height,
		DisableResize:     true,
		Frameless:         true,
		BackgroundColour:  application.NewRGB(10, 10, 16),
		URL:               "/",
		MaximiseButtonState: application.ButtonHidden,
	})

	launcher.SetWindow(window)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
