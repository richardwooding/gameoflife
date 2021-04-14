package life

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"io"
)

type Life struct {
	app.Compo
	generation int64
	dx int
	dy int
	colony *[][]bool
}

func New(dx uint, dy uint) *Life {
	c := make([][]bool, dy)
	for i := range c {
		c[i] = make([]bool, dx)
	}
	return &Life{
		generation: 0,
		dx:   int(dx),
		dy:   int(dy),
		colony: &c,
	}
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

func  (l *Life) countNeighbours(x int, y int) uint {
	return l.count(x-1, y-1) + l.count(x, y-1) + l.count(x+1, y-1) +
		   l.count(x-1, y)                     + l.count(x+1, y)   +
		   l.count(x-1, y+1) + l.count(x, y+1) + l.count(x+1, y+1)
}

func (l *Life) Generate() {
	ng := make([][]bool, l.dy)
	for i := range ng {
		ng[i] = make([]bool, l.dx)
	}
	for x := 0; x<l.dx; x++ {
		for y := 0; y<l.dy; y++ {
			alive := (*l.colony)[y][x]
			neighbours := l.countNeighbours(x, y)
			ng[y][x] =  (alive && (neighbours == 2 || neighbours == 3)) || (!alive && neighbours == 3)
		}
	}
	l.generation++
	l.colony = &ng
}

func (l *Life) Alive(x int, y int) {
	(*l.colony)[y][x] = true
}

func (l *Life) Dead(x int, y int) {
	(*l.colony)[y][x] = false
}

func (l *Life) Toggle(x int, y int) {
	(*l.colony)[y][x] = !(*l.colony)[y][x]
}

func (l *Life) Print(w io.Writer) {
  fmt.Fprintf(w, "Gneration: %d\n", l.generation)
  for y := range *l.colony {
  	for x := range (*l.colony)[y] {
  		if (*l.colony)[y][x] {
  			fmt.Fprint(w, "*")
		} else {
			fmt.Fprint(w, "-")
		}
	}
	fmt.Fprintln(w)
  }
}

func (l *Life) GenerateAndPrint(w io.Writer) {
	l.Generate()
	l.Print(w)
}