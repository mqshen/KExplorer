package main

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"kafkaexplorer/backend/services"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var version = "0.0.0"

func main() {
	// Create an instance of the app structure
	sysSvc := services.System()
	connSvc := services.Connection()
	topicSvc := services.Topic()
	browserSvc := services.Browser()
	prefSvc := services.Preferences()

	prefSvc.SetAppVersion(version)
	windowWidth, windowHeight, maximised := prefSvc.GetWindowSize()
	windowStartState := options.Normal
	if maximised {
		windowStartState = options.Maximised
	}

	// menu
	appMenu := menu.NewMenu()
	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.AppMenu())
		appMenu.Append(menu.EditMenu())
		appMenu.Append(menu.WindowMenu())
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "KExplorer",
		Width:            windowWidth,
		Height:           windowHeight,
		WindowStartState: windowStartState,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        runtime.GOOS != "darwin",
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			sysSvc.Start(ctx, version)
			connSvc.Start(ctx)
			topicSvc.Start(ctx)
			browserSvc.Start(ctx)
			//cliSvc.Start(ctx)
			//monitorSvc.Start(ctx)
			//pubsubSvc.Start(ctx)
			//
			//services.GA().SetSecretKey(gaMeasurementID, gaSecretKey)
			//services.GA().Startup(version)
		},
		OnDomReady: func(ctx context.Context) {
			x, y := prefSvc.GetWindowPosition(ctx)
			wailsRuntime.WindowSetPosition(ctx, x, y)
			wailsRuntime.WindowShow(ctx)
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			x, y := wailsRuntime.WindowGetPosition(ctx)
			prefSvc.SaveWindowPosition(x, y)
			return false
		},
		OnShutdown: func(ctx context.Context) {
			browserSvc.Stop()
			//cliSvc.CloseAll()
			//monitorSvc.StopAll()
			//pubsubSvc.StopAll()
		},
		Bind: []interface{}{
			sysSvc,
			connSvc,
			topicSvc,
			browserSvc,
			//cliSvc,
			//monitorSvc,
			//pubsubSvc,
			prefSvc,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "Kafka Explorer" + version,
				Message: "A modern lightweight cross-platform Kafka desktop client.\n\nCopyright © 2024",
				Icon:    icon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
