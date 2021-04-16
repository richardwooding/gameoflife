package life

import (
	"github.com/enescakir/emoji"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"time"
)

type Life struct {
	app.Compo
	generation int64
	dx         int
	dy         int
	colony     *[][]bool
	ticker     *time.Ticker
	done	   chan bool
}

func (l *Life) newColony(dx uint, dy uint) {
	c := make([][]bool, dy)
	for i := range c {
		c[i] = make([]bool, dx)
	}
	l.generation = 0
	l.dx = int(dx)
	l.dy = int(dy)
	l.colony = &c
	l.Update()
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

func (l *Life) generate() {
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
	l.Update()
}

func (l *Life) alive(x int, y int) {
	(*l.colony)[y][x] = true
	l.Update()
}

func (l *Life) dead(x int, y int) {
	(*l.colony)[y][x] = false
	l.Update()
}

func (l *Life) toggle(x int, y int) {
	(*l.colony)[y][x] = !(*l.colony)[y][x]
	l.Update()
}

func (l *Life) className(x int, y int) string {
	if (*l.colony)[y][x] {
		return "alive"
	} else {
		return "dead"
	}
}

func (l *Life) startTicking(ctx app.Context) {
	l.ticker = time.NewTicker(50*time.Millisecond)
	l.done = make(chan bool)
	ctx.Async(func() {
		for {
			select {
			case <-l.done:
				return
			case <- l.ticker.C:
				l.Defer(func(context app.Context) {
					l.generate()
				})
			}
		}
	})
}

func (l *Life) stopTicking() {
	l.ticker.Stop()
	l.ticker = nil
	l.Update()
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
			app.Button().Text("Make Colony").OnClick(func(ctx app.Context, e app.Event) {
				l.newColony(64, 64)
			}),
		).Else(
			app.If(l.ticker == nil,
				app.Button().Text(emoji.PlayButton).OnClick(func(ctx app.Context, e app.Event) {
					l.startTicking(ctx)
				}),
			).Else(
				app.Button().Text(emoji.PauseButton).OnClick(func(ctx app.Context, e app.Event) {
					l.stopTicking()
				}),
			),
			app.Hr(),
			app.Div().Class("wrapper").Body(
				app.Range(colony).Slice(func(y int) app.UI {
					return app.Range(colony[y]).Slice(func(x int) app.UI {
						return app.Div().Class(l.className(x, y)).OnClick(func(ctx app.Context, e app.Event) {
							if l.ticker == nil {
								l.toggle(x, y)
							}
						})
					})
				}),
			),
		),
	)
}
