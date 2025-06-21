package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/richardwooding/gameoflife/pkg/life"
	"github.com/richardwooding/gameoflife/webmode"
	"log"
	"net/http"
	"os"
)

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {

	// Components routing:
	app.RouteWithRegexp("/(.*)", func() app.Composer { return &life.Life{} })
	app.RunWhenOnBrowser()

	webMode := webmode.Live
	webModeEnv, ok := os.LookupEnv("CONWAYS_GAME_OF_LIFE_WEB_MODE")
	if ok {
		if parsedWebMode, err := webmode.ParseWebMode(webModeEnv); err == nil {
			webMode = parsedWebMode
		} else {
			log.Fatal(err)
		}
	}

	handler := &app.Handler{
		Name:        "Conway's Game of Life",
		Description: "A live demo of Conway's Game of Life.",
		Styles: []string{
			"/web/gameoflife.css",
		},
	}

	switch webMode {
	case webmode.Live:
		// HTTP routing:
		http.Handle("/{path...}", handler)

		if err := http.ListenAndServe(":8000", nil); err != nil {
			log.Fatal(err)
		}
	case webmode.Static:
		if err := app.GenerateStaticWebsite(".", handler); err != nil {
			log.Fatal(err)
		}
	}

}
