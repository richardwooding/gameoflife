package main

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/richardwooding/gameoflife/pkg/life"
	"log"
	"net/http"
)


// The main function is the entry point where the app is configured and started.
// It is executed in 2 different environments: A client (the web browser) and a
// server.
func main() {

	// The first thing to do is to associate the hello component with a path.
	//
	// This is done by calling the Route() function,  which tells go-app what
	// component to display for a given path, on both client and server-side.
	app.RouteWithRegexp("^/.*", &life.Life{})

	// Once the routes set up, the next thing to do is to either launch the app
	// or the server that serves the app.
	//
	// When executed on the client-side, the RunWhenOnBrowser() function
	// launches the app,  starting a loop that listens for app events and
	// executes client instructions. Since it is a blocking call, the code below
	// it will never be executed.
	//
	// On the server-side, RunWhenOnBrowser() does nothing, which allows the
	// writing of server logic without needing precompiling instructions.
	app.RunWhenOnBrowser()


		err := app.GenerateStaticWebsite("docs", &app.Handler{
			Name:        "Conway's Game Of Life",
			Description: "Conway's Game Of Lifer",
			Styles: []string {
				"/web/app.css",
			},
			Resources: app.GitHubPages("gameoflife"),
		})

		if err != nil {
			log.Fatal(err)
		}


	// Finally, launching the server that serves the app is done by using the Go
	// standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client and all its
	// required resources to make it work into a web browser. Here it is
	// configured to handle requests with a path that starts with "/".
	http.Handle("/", &app.Handler{
		Name:        "Conway's Game Of Life",
		Description: "Conway's Game Of Life",
		Styles: []string {
			"/web/app.css",
		},
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
