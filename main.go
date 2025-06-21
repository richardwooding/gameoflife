package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/richardwooding/gameoflife/pkg/life"
	"log"
	"net/http"
)

// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {

	// Components routing:
	app.RouteWithRegexp("/(.*)", func() app.Composer { return &life.Life{} })
	app.RunWhenOnBrowser()

	// HTTP routing:
	http.Handle("/{path...}", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
		Styles: []string{
			"/web/gameoflife.css",
		},
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
