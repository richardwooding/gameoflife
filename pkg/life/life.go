package life

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Life struct {
	app.Compo
	generation   int64
	dx           int
	dy           int
	colony       *[][]bool
	ticker       *time.Ticker
	done         chan bool
	tickInterval time.Duration
}

type exported struct {
	Colony [][]bool
}

// newColony initializes a new colony with the given dimensions and resets the simulation state.
func (l *Life) newColony(context app.Context, dx uint, dy uint) {
	c := make([][]bool, dy)
	for i := range c {
		c[i] = make([]bool, dx)
	}
	l.generation = 0
	l.dx = int(dx)
	l.dy = int(dy)
	l.colony = &c
	l.tickInterval = 50 * time.Millisecond
	l.saveState(context)
}

// count returns 1 if the cell at (x, y) is alive, 0 otherwise.
func (l *Life) count(x int, y int) uint {
	if x < 0 || y < 0 || x >= l.dx || y >= l.dy {
		return 0
	}
	if (*l.colony)[y][x] {
		return 1
	}
	return 0
}

// countNeighbours counts the alive neighbours of the cell at (x, y).
func (l *Life) countNeighbours(x int, y int) uint {
	return l.count(x-1, y-1) + l.count(x, y-1) + l.count(x+1, y-1) +
		l.count(x-1, y) + l.count(x+1, y) +
		l.count(x-1, y+1) + l.count(x, y+1) + l.count(x+1, y+1)
}

// generate computes the next generation of the colony based on the current state.
func (l *Life) generate(ctx app.Context) {
	ng := make([][]bool, l.dy)
	for i := range ng {
		ng[i] = make([]bool, l.dx)
	}
	for x := 0; x < l.dx; x++ {
		for y := 0; y < l.dy; y++ {
			alive := (*l.colony)[y][x]
			neighbours := l.countNeighbours(x, y)
			ng[y][x] = (alive && (neighbours == 2 || neighbours == 3)) || (!alive && neighbours == 3)
		}
	}
	l.generation++
	l.colony = &ng
	ctx.Update()
}

// toggle toggles the alive state of the cell at (x, y) and saves the current state.
func (l *Life) toggle(context app.Context, x int, y int) {
	(*l.colony)[y][x] = !(*l.colony)[y][x]
	l.saveState(context)
	//l.Update()
}

// className returns the CSS class name ("alive" or "dead") for the cell at (x, y).
func (l *Life) className(x int, y int) string {
	if (*l.colony)[y][x] {
		return "alive"
	} else {
		return "dead"
	}
}

// startTicking starts the simulation ticker with the current tick interval.
func (l *Life) startTicking(ctx app.Context) {
	if l.tickInterval == 0 {
		l.tickInterval = 50 * time.Millisecond
	}
	l.ticker = time.NewTicker(l.tickInterval)
	l.done = make(chan bool)
	ctx.Async(func() {
		for {
			select {
			case <-l.done:
				return
			case <-l.ticker.C:
				l.generate(ctx)
			}
		}
	})
}

// stopTicking stops the simulation ticker and saves the current state.
func (l *Life) stopTicking(ctx app.Context) {
	l.ticker.Stop()
	l.ticker = nil
	l.saveState(ctx)
}

// OnNav loads the simulation state from the URL path if present.
func (l *Life) OnNav(ctx app.Context) {
	path := ctx.Page().URL().Path
	if path != "/" && path != "/gameoflife/" {
		l.loadState(strings.TrimPrefix(path, "/gameoflife")[1:])
	}
}

// OnMount is called when the component is mounted. (Implementation may be added as needed.)
func (l *Life) OnMount(ctx app.Context) {
	fragment := ctx.Page().URL().Fragment
	if fragment != "" {
		l.loadState(fragment)
	}
}

// loadState decodes and loads the simulation state from a base64-encoded string.
func (l *Life) loadState(state string) {
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
		l.tickInterval = 50 * time.Millisecond
		l.dy = len((*exp).Colony)
		l.dx = len((*exp).Colony[0])
		l.colony = &(*exp).Colony
		//l.Update()
	}
}

// saveState encodes and saves the current simulation state as a base64-encoded string in the URL.
func (l *Life) saveState(context app.Context) {
	exp := exported{Colony: *l.colony}
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
func (l *Life) clearColony(ctx app.Context) {
	if l.colony == nil {
		return
	}
	for y := 0; y < l.dy; y++ {
		for x := 0; x < l.dx; x++ {
			(*l.colony)[y][x] = false
		}
	}
	l.generation = 0
	l.saveState(ctx)
}

// insertGlider inserts a glider pattern at the top-left corner of the colony.
func (l *Life) insertGlider(ctx app.Context) {
	if l.colony == nil || l.dx < 3 || l.dy < 3 {
		return
	}
	// Clear a 3x3 area at the top-left
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			(*l.colony)[y][x] = false
		}
	}
	// Set glider pattern
	(*l.colony)[0][1] = true
	(*l.colony)[1][2] = true
	(*l.colony)[2][0] = true
	(*l.colony)[2][1] = true
	(*l.colony)[2][2] = true
	l.saveState(ctx)
	ctx.Update()
}

// insertBlinker inserts a blinker pattern at the center of the colony.
func (l *Life) insertBlinker(ctx app.Context) {
	if l.colony == nil || l.dx < 5 || l.dy < 5 {
		return
	}
	// Place blinker at (2,2) horizontally, clear a 5x1 area at y=2
	for i := 0; i < 5; i++ {
		(*l.colony)[2][i] = false
	}
	// Set blinker pattern at (2,2)-(4,2)
	(*l.colony)[2][2] = true
	(*l.colony)[2][3] = true
	(*l.colony)[2][4] = true
	l.saveState(ctx)
	ctx.Update()
}

// insertPulsar inserts a pulsar pattern at the specified location in the colony.
func (l *Life) insertPulsar(ctx app.Context) {
	if l.colony == nil || l.dx < 9 || l.dy < 9 {
		return
	}
	// Place pulsar at (2,2), clear a 7x7 area at (2,2)
	for y := 2; y < 9; y++ {
		for x := 2; x < 9; x++ {
			(*l.colony)[y][x] = false
		}
	}
	// Set pulsar pattern at (2,2)
	(*l.colony)[3][7] = true
	(*l.colony)[4][7] = true
	(*l.colony)[5][7] = true
	(*l.colony)[6][7] = true
	(*l.colony)[7][7] = true
	(*l.colony)[8][7] = true
	(*l.colony)[7][3] = true
	(*l.colony)[7][4] = true
	(*l.colony)[7][5] = true
	(*l.colony)[7][6] = true
	(*l.colony)[7][8] = true
	(*l.colony)[7][9] = true
	l.saveState(ctx)
	ctx.Update()
}

// insertToad inserts a toad pattern at the specified location in the colony.
func (l *Life) insertToad(ctx app.Context) {
	if l.colony == nil || l.dx < 6 || l.dy < 6 {
		return
	}
	for y := 2; y < 5; y++ {
		for x := 2; x < 6; x++ {
			(*l.colony)[y][x] = false
		}
	}
	(*l.colony)[2][3] = true
	(*l.colony)[2][4] = true
	(*l.colony)[2][5] = true
	(*l.colony)[3][2] = true
	(*l.colony)[3][3] = true
	(*l.colony)[3][4] = true
	l.saveState(ctx)
	ctx.Update()
}

// insertBeacon inserts a beacon pattern at the specified location in the colony.
func (l *Life) insertBeacon(ctx app.Context) {
	if l.colony == nil || l.dx < 6 || l.dy < 6 {
		return
	}
	for y := 2; y < 6; y++ {
		for x := 2; x < 6; x++ {
			(*l.colony)[y][x] = false
		}
	}
	(*l.colony)[2][2] = true
	(*l.colony)[2][3] = true
	(*l.colony)[3][2] = true
	(*l.colony)[3][3] = true
	(*l.colony)[4][4] = true
	(*l.colony)[4][5] = true
	(*l.colony)[5][4] = true
	(*l.colony)[5][5] = true
	l.saveState(ctx)
	ctx.Update()
}

// insertAcorn inserts an acorn pattern at the specified location in the colony.
func (l *Life) insertAcorn(ctx app.Context) {
	if l.colony == nil || l.dx < 10 || l.dy < 10 {
		return
	}
	for y := 2; y < 5; y++ {
		for x := 2; x < 9; x++ {
			(*l.colony)[y][x] = false
		}
	}
	(*l.colony)[3][3] = true
	(*l.colony)[4][5] = true
	(*l.colony)[2][4] = true
	(*l.colony)[3][5] = true
	(*l.colony)[3][6] = true
	(*l.colony)[3][7] = true
	(*l.colony)[3][8] = true
	l.saveState(ctx)
	ctx.Update()
}

// insertGosperGliderGun inserts a Gosper Glider Gun pattern at the specified location in the colony.
func (l *Life) insertGosperGliderGun(ctx app.Context) {
	if l.colony == nil || l.dx < 40 || l.dy < 11 {
		return
	}
	for y := 2; y < 11; y++ {
		for x := 2; x < 39; x++ {
			(*l.colony)[y][x] = false
		}
	}
	// Set Gosper Glider Gun pattern (relative to (2,2))
	cells := [][2]int{
		{2, 6}, {2, 7}, {3, 6}, {3, 7},
		{12, 6}, {12, 7}, {12, 8}, {13, 5}, {13, 9}, {14, 4}, {14, 10}, {15, 4}, {15, 10}, {16, 7},
		{17, 5}, {17, 9}, {18, 6}, {18, 7}, {18, 8}, {19, 7},
		{22, 4}, {22, 5}, {22, 6}, {23, 4}, {23, 5}, {23, 6}, {24, 3}, {24, 7}, {26, 2}, {26, 3}, {26, 7}, {26, 8},
		{36, 4}, {36, 5}, {37, 4}, {37, 5},
	}
	for _, c := range cells {
		(*l.colony)[c[1]][c[0]] = true
	}
	l.saveState(ctx)
	ctx.Update()
}

// setSpeed adjusts the simulation speed and restarts the ticker with the new interval.
func (l *Life) setSpeed(ctx app.Context, ms int64) {
	if ms < 10 {
		ms = 10
	}
	if ms > 1000 {
		ms = 1000
	}
	l.tickInterval = time.Duration(ms) * time.Millisecond
	if l.ticker != nil {
		l.stopTicking(ctx)
		l.startTicking(ctx)
	}
}

// Render generates the UI for the Game of Life component.
func (l *Life) Render() app.UI {

	var colony [][]bool
	if l.colony == nil {
		colony = [][]bool{}
	} else {
		colony = *l.colony
	}

	return app.Div().Body(
		app.H1().Text("Game of life"),
		app.If(l.colony == nil,
			func() app.UI {
				return app.Button().Text("Make Colony").OnClick(func(ctx app.Context, e app.Event) {
					l.newColony(ctx, 64, 64)
					l.tickInterval = 50 * time.Millisecond
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
						Value(fmt.Sprintf("%d", l.tickInterval.Milliseconds())).
						Aria("label", "Simulation speed interval in milliseconds").
						Aria("valuenow", fmt.Sprintf("%d", l.tickInterval.Milliseconds())).
						Aria("valuemin", "10").
						Aria("valuemax", "1000").
						OnInput(func(ctx app.Context, e app.Event) {
							if targetSpeed, err := strconv.Atoi(e.Get("target").Get("value").String()); err == nil {
								l.setSpeed(ctx, int64(targetSpeed))
							}
						}),
					app.Span().Style("margin-left", "8px").Textf("%d ms", l.tickInterval.Milliseconds()),
				),
				// Play/Pause and other controls
				app.If(l.ticker == nil,
					func() app.UI {
						return app.Button().Text(emoji.PlayButton).OnClick(func(ctx app.Context, e app.Event) {
							l.startTicking(ctx)
						})
					}).Else(func() app.UI {
					return app.Button().Text(emoji.PauseButton).OnClick(func(ctx app.Context, e app.Event) {
						l.stopTicking(ctx)
					})
				}),
				app.Button().Text(emoji.ClButton).OnClick(func(ctx app.Context, e app.Event) {
					if l.ticker == nil {
						l.clearColony(ctx)
					}
				}),
				app.Button().Textf("%s Glider", emoji.Plus).OnClick(func(ctx app.Context, e app.Event) {
					if l.ticker == nil {
						l.insertGlider(ctx)
					}
				}),
				app.Button().Textf("%s Blinker", emoji.Plus).OnClick(func(ctx app.Context, e app.Event) {
					if l.ticker == nil {
						l.insertBlinker(ctx)
					}
				}),
				app.Button().Textf("%s Pulsar", emoji.Plus).OnClick(func(ctx app.Context, e app.Event) {
					if l.ticker == nil {
						l.insertPulsar(ctx)
					}
				}),
				app.Button().Textf("%s Toad", emoji.Plus).OnClick(func(ctx app.Context, e app.Event) {
					if l.ticker == nil {
						l.insertToad(ctx)
					}
				}),
				app.Button().Textf("%s Beacon", emoji.Plus).OnClick(func(ctx app.Context, e app.Event) {
					if l.ticker == nil {
						l.insertBeacon(ctx)
					}
				}),
				app.Button().Textf("%s Acorn", emoji.Plus).OnClick(func(ctx app.Context, e app.Event) {
					if l.ticker == nil {
						l.insertAcorn(ctx)
					}
				}),
				app.Button().Textf("%s Gosper Glider Gun", emoji.Plus).OnClick(func(ctx app.Context, e app.Event) {
					if l.ticker == nil {
						l.insertGosperGliderGun(ctx)
					}
				}),
			)
		}),
		app.Hr(),
		app.Div().Textf("Generation: %d", l.generation),
		app.Div().Class("wrapper").Body(
			app.Range(colony).Slice(func(y int) app.UI {
				return app.Range(colony[y]).Slice(func(x int) app.UI {
					return app.Div().Class(l.className(x, y)).OnClick(func(ctx app.Context, e app.Event) {
						if l.ticker == nil {
							l.toggle(ctx, x, y)
						}
					})
				})
			}),
		),
	)
}
