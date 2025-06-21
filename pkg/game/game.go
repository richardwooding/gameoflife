package game

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/richardwooding/gameoflife/model"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	app.Compo
	colony       *model.Colony
	ticker       *time.Ticker
	done         chan bool
	tickInterval time.Duration
}

type exported struct {
	Cells [][]bool
}

// NewColony initializes a new colony with the given dimensions and resets the simulation state.
func (g *Game) NewColony(context app.Context, dx uint, dy uint) {
	g.colony = model.NewColony(int(dx), int(dy))
	g.tickInterval = 50 * time.Millisecond
	g.saveState(context)
}

func (g *Game) Generate(ctx app.Context) {
	g.colony.Generate()
	ctx.Update()
}

// toggle toggles the alive state of the cell at (x, y) and saves the current state.
func (g *Game) toggle(context app.Context, x int, y int) {
	g.colony.Toggle(x, y)
	g.saveState(context)
}

// className returns the CSS class name ("alive" or "dead") for the cell at (x, y).
func (g *Game) className(x int, y int) string {
	if g.colony.IsAlive(x, y) {
		return "alive"
	} else {
		return "dead"
	}
}

// startTicking starts the simulation ticker with the current tick intervag.
func (g *Game) startTicking(ctx app.Context) {
	if g.tickInterval == 0 {
		g.tickInterval = 50 * time.Millisecond
	}
	g.ticker = time.NewTicker(g.tickInterval)
	g.done = make(chan bool)
	ctx.Async(func() {
		for {
			select {
			case <-g.done:
				return
			case <-g.ticker.C:
				g.Generate(ctx)
			}
		}
	})
}

// stopTicking stops the simulation ticker and saves the current state.
func (g *Game) stopTicking(ctx app.Context) {
	g.ticker.Stop()
	g.ticker = nil
	g.saveState(ctx)
}

// OnNav loads the simulation state from the URL path if present.
func (g *Game) OnNav(ctx app.Context) {
	path := ctx.Page().URL().Path
	if path != "/" && path != "/gameoflife/" {
		g.loadState(strings.TrimPrefix(path, "/gameoflife")[1:])
	}
}

// OnMount is called when the component is mounted. (Implementation may be added as needed.)
func (g *Game) OnMount(ctx app.Context) {
	fragment := ctx.Page().URL().Fragment
	if fragment != "" {
		g.loadState(fragment)
	}
}

// loadState decodes and loads the simulation state from a base64-encoded string.
func (g *Game) loadState(state string) {
	exp := &exported{}
	b, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return
	}
	buff := bytes.NewBuffer(b)
	reader := flate.NewReader(buff)

	dec := gob.NewDecoder(reader)
	err = dec.Decode(exp)
	if err == nil {
		g.tickInterval = 50 * time.Millisecond
		if g.colony == nil {
			g.colony = model.NewColony(len(exp.Cells[0]), len(exp.Cells))
		} else {
			g.colony.Reset()
		}
		g.colony.SetCells(exp.Cells)
		//g.Update()
	}
}

// saveState encodes and saves the current simulation state as a base64-encoded string in the URL.
func (g *Game) saveState(context app.Context) {
	exp := exported{Cells: *g.colony.Cells()}
	var buff bytes.Buffer
	writer, _ := flate.NewWriter(&buff, flate.BestCompression)
	enc := gob.NewEncoder(writer)
	_ = enc.Encode(exp)
	_ = writer.Flush()
	str := base64.RawURLEncoding.EncodeToString(buff.Bytes())
	path := context.Page().URL().Path
	var prefix string
	if strings.HasPrefix(path, "/gameoflife") {
		prefix = "/gameoflife/"
	} else {
		prefix = "/"
	}
	newUrl, _ := url.Parse(fmt.Sprintf("%s#%s", prefix, str))
	context.Page().ReplaceURL(context.Page().URL().ResolveReference(newUrl))
}

// clearColony clears the colony, resetting all cells to dead.
func (g *Game) clearColony(ctx app.Context) {
	if g.colony == nil {
		return
	}
	g.colony.Reset()
	g.saveState(ctx)
}

// setSpeed adjusts the simulation speed and restarts the ticker with the new intervag.
func (g *Game) setSpeed(ctx app.Context, ms int64) {
	if ms < 10 {
		ms = 10
	}
	if ms > 1000 {
		ms = 1000
	}
	g.tickInterval = time.Duration(ms) * time.Millisecond
	if g.ticker != nil {
		g.stopTicking(ctx)
		g.startTicking(ctx)
	}
}

// Render generates the UI for the Game of Game component.
func (g *Game) Render() app.UI {

	return app.Div().Body(
		app.H1().Text("Conway's Game of life"),
		app.Button().Textf("%s Open on Github", emoji.Laptop).OnClick(func(ctx app.Context, e app.Event) {
			ctx.Navigate("https://github.com/richardwooding/gameoflife")
		}),
		app.If(g.colony == nil,
			func() app.UI {
				return app.Button().Textf("%s Make Cells", emoji.Hut).OnClick(func(ctx app.Context, e app.Event) {
					g.NewColony(ctx, 64, 64)
					g.tickInterval = 50 * time.Millisecond
				})
			}).Else(func() app.UI {
			return app.Div().Body(
				// Range slider for speed
				app.Div().Body(
					app.Label().Text("Interval: ").For("interval-slider"),
					app.Input().
						Type("range").
						ID("interval-slider").
						Min("10").
						Max("1000").
						Step(10).
						Value(fmt.Sprintf("%d", g.tickInterval.Milliseconds())).
						Aria("label", "Simulation speed interval in milliseconds").
						Aria("valuenow", fmt.Sprintf("%d", g.tickInterval.Milliseconds())).
						Aria("valuemin", "10").
						Aria("valuemax", "1000").
						OnInput(func(ctx app.Context, e app.Event) {
							if targetSpeed, err := strconv.Atoi(e.Get("target").Get("value").String()); err == nil {
								g.setSpeed(ctx, int64(targetSpeed))
							}
						}),
					app.Span().Style("margin-left", "8px").Textf("%d ms", g.tickInterval.Milliseconds()),
				),
				// Play/Pause and other controls
				app.If(g.ticker == nil,
					func() app.UI {
						return app.Button().Text(emoji.PlayButton).OnClick(func(ctx app.Context, e app.Event) {
							g.startTicking(ctx)
						})
					}).Else(func() app.UI {
					return app.Button().Text(emoji.PauseButton).OnClick(func(ctx app.Context, e app.Event) {
						g.stopTicking(ctx)
					})
				}),
				app.Button().Text(emoji.ClButton).OnClick(func(ctx app.Context, e app.Event) {
					if g.ticker == nil {
						g.clearColony(ctx)
					}
				}),
				app.Range(Patterns).Slice(func(i int) app.UI {
					return app.Button().Textf("%s %s", emoji.Plus, Patterns[i].GetName()).OnClick(func(ctx app.Context, e app.Event) {
						Patterns[i].Stamp(g.colony.Cells(), 2, 2)
						g.saveState(ctx)
					})
				}),
				app.Button().Textf("%s Random", emoji.GameDie).OnClick(func(ctx app.Context, e app.Event) {
					if g.colony != nil && g.ticker == nil {
						g.insertRandom(ctx)
					}
				}),
				app.Button().Textf("%s Center", emoji.Compass).OnClick(func(ctx app.Context, e app.Event) {
					if g.colony != nil && g.ticker == nil {
						g.centerAlive(ctx)
					}
				}),
			)
		}),
		app.Hr(),
		app.If(g.colony != nil, func() app.UI {
			return app.Div().Textf("Generation: %d", g.colony.GetGeneration())
		}),
		app.If(g.colony != nil, func() app.UI {
			return app.Div().Class("wrapper").Body(
				app.Range(*g.colony.Cells()).Slice(func(y int) app.UI {
					return app.Range((*g.colony.Cells())[y]).Slice(func(x int) app.UI {
						return app.Div().Class(g.className(x, y)).OnClick(func(ctx app.Context, e app.Event) {
							if g.ticker == nil {
								g.toggle(ctx, x, y)
							}
						})
					})
				}))
		}),
	)
}

func (g *Game) insertRandom(ctx app.Context) {
	g.colony.Randomize()
	g.saveState(ctx)
}

// centerAlive shifts the bounding box of alive cells to the center of the grid.
func (g *Game) centerAlive(ctx app.Context) {
	g.colony.CentreAlive()
	g.saveState(ctx)
	ctx.Update()
}
