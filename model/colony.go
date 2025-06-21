package model

import "math/rand"

type Colony struct {
	generation int64
	dx         int
	dy         int
	cells      *[][]bool
}

func NewColony(dx, dy int) *Colony {
	cells := make([][]bool, dy)
	for i := range cells {
		cells[i] = make([]bool, dx)
	}
	return &Colony{
		generation: 0,
		dx:         dx,
		dy:         dy,
		cells:      &cells,
	}
}

func (c *Colony) SetCells(cells [][]bool) {
	c.dy = len(cells)
	c.dy = len(cells[0])
	c.cells = &cells
}

func (c *Colony) Cells() *[][]bool {
	return c.cells
}

// count returns 1 if the cell at (x, y) is alive, 0 otherwise.
func (c *Colony) count(x int, y int) uint {
	if x < 0 || y < 0 || x >= c.dx || y >= c.dy {
		return 0
	}
	if (*c.cells)[y][x] {
		return 1
	}
	return 0
}

// countNeighbours counts the alive neighbours of the cell at (x, y).
func (c *Colony) countNeighbours(x int, y int) uint {
	return c.count(x-1, y-1) + c.count(x, y-1) + c.count(x+1, y-1) +
		c.count(x-1, y) + c.count(x+1, y) +
		c.count(x-1, y+1) + c.count(x, y+1) + c.count(x+1, y+1)
}

func (c *Colony) GetGeneration() int64 {
	if c == nil {
		return 0
	}
	return c.generation
}

func (c *Colony) Generate() {
	ng := make([][]bool, c.dy)
	for i := range ng {
		ng[i] = make([]bool, c.dx)
	}
	for x := 0; x < c.dx; x++ {
		for y := 0; y < c.dy; y++ {
			alive := (*c.cells)[y][x]
			neighbours := c.countNeighbours(x, y)
			ng[y][x] = (alive && (neighbours == 2 || neighbours == 3)) || (!alive && neighbours == 3)
		}
	}
	c.cells = &ng
	c.generation++
}

func (c *Colony) Toggle(x, y int) {
	if x < 0 || y < 0 || x >= c.dx || y >= c.dy {
		return
	}
	(*c.cells)[y][x] = !(*c.cells)[y][x]
}

func (c *Colony) IsAlive(x, y int) bool {
	if x < 0 || y < 0 || x >= c.dx || y >= c.dy {
		return false
	}
	return (*c.cells)[y][x]
}

func (c *Colony) Reset() {
	c.generation = 0
	for y := 0; y < c.dy; y++ {
		for x := 0; x < c.dx; x++ {
			(*c.cells)[y][x] = false
		}
	}
}

func (c *Colony) Randomize() {
	for y := 0; y < c.dy; y++ {
		for x := 0; x < c.dx; x++ {
			(*c.cells)[y][x] = rand.Intn(2) == 1
		}
	}
	c.generation = 0
}

func (c *Colony) CentreAlive() {
	minX, minY, maxX, maxY := c.dx, c.dy, 0, 0
	found := false
	for y := 0; y < c.dy; y++ {
		for x := 0; x < c.dx; x++ {
			if (*c.cells)[y][x] {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
				found = true
			}
		}
	}
	if !found {
		return
	}
	w, h := maxX-minX+1, maxY-minY+1
	dx := (c.dx-w)/2 - minX
	dy := (c.dy-h)/2 - minY
	if dx == 0 && dy == 0 {
		return
	}
	newCells := make([][]bool, c.dy)
	for y := range newCells {
		newCells[y] = make([]bool, c.dx)
	}
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if (*c.cells)[y][x] {
				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < c.dx && ny >= 0 && ny < c.dy {
					newCells[ny][nx] = true
				}
			}
		}
	}
	c.cells = &newCells
}
