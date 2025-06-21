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

func (l *Life) count(x int, y int) uint {
	if x < 0 || y < 0 || x >= l.dx || y >= l.dy {
		return 0
	}
	if (*l.colony)[y][x] {
		return 1
	}
	return 0
}

func (l *Life) countNeighbours(x int, y int) uint {
	return l.count(x-1, y-1) + l.count(x, y-1) + l.count(x+1, y-1) +
		l.count(x-1, y) + l.count(x+1, y) +
		l.count(x-1, y+1) + l.count(x, y+1) + l.count(x+1, y+1)
}

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

func (l *Life) toggle(context app.Context, x int, y int) {
	(*l.colony)[y][x] = !(*l.colony)[y][x]
	l.saveState(context)
	//l.Update()
}

func (l *Life) className(x int, y int) string {
	if (*l.colony)[y][x] {
		return "alive"
	} else {
		return "dead"
	}
}

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

func (l *Life) stopTicking(ctx app.Context) {
	l.ticker.Stop()
	l.ticker = nil
	l.saveState(ctx)
}

func (l *Life) OnNav(ctx app.Context) {
	path := ctx.Page().URL().Path
	if path != "/" && path != "/gameoflife/" {
		l.loadState(strings.TrimPrefix(path, "/gameoflife")[1:])
	}
}

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

func (l *Life) saveState(context app.Context) {
	exp := exported{Colony: *l.colony}
	var buff bytes.Buffer
	writer, _ := flate.NewWriter(&buff, flate.BestCompression)
	enc := gob.NewEncoder(writer)
	_ = enc.Encode(exp)
	_ = writer.Flush()
	str := base64.RawURLEncoding.EncodeToString(buff.Bytes())
	println("Saving state", str)
	path := context.Page().URL().Path
	var prefix string
	if strings.HasPrefix(path, "/gameoflife") {
		prefix = "/gameoflife/"
	} else {
		prefix = "/"
	}
	newUrl, _ := url.Parse(fmt.Sprintf("%s%s", prefix, str))
	context.Page().ReplaceURL(context.Page().URL().ResolveReference(newUrl))
}

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

func (l *Life) setSpeed(ctx app.Context, ms int64) {
	println("Setting speed to", ms, "ms")
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
