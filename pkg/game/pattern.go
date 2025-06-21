package game

// Pattern represents a named pattern using a sparse list of live cell coordinates.
type Pattern struct {
	name   string   // Name of the feature (pattern)
	sparse [][2]int // List of live cell coordinates (x, y)
}

// GetName returns the name of the feature.
func (p *Pattern) GetName() string {
	return p.name
}

// Stamp applies the feature's live cells to the given grid at the specified offset.
func (p *Pattern) Stamp(cells *[][]bool, offsetX, offsetY int) {
	for _, pos := range p.sparse {
		x, y := pos[0]+offsetX, pos[1]+offsetY
		// Ensure we don't go out of bounds
		if y >= 0 && y < len(*cells) && x >= 0 && x < len((*cells)[0]) {
			(*cells)[y][x] = true
		}
	}
}

// Patterns is a list of predefined Game of Life patterns, each using sparse representation.
var Patterns = []Pattern{
	// Glider: A small pattern that moves diagonally across the grid.
	{
		name: "Glider",
		sparse: [][2]int{
			{1, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2},
		},
	},
	// Blinker: A simple oscillator with period 2 (flips between vertical and horizontal line).
	{
		name: "Blinker",
		sparse: [][2]int{
			{0, 1}, {1, 1}, {2, 1},
		},
	},
	// Pulsar: A larger period-3 oscillator with a distinctive shape.
	{
		name: "Pulsar",
		sparse: [][2]int{
			{3, 0}, {4, 0}, {5, 0}, {7, 0}, {8, 0}, {9, 0},
			{2, 1}, {4, 1}, {6, 1}, {8, 1}, {10, 1},
			{0, 3}, {1, 3}, {4, 3}, {5, 3}, {6, 3}, {7, 3}, {8, 3}, {11, 3}, {12, 3},
			{0, 4}, {2, 4}, {4, 4}, {6, 4}, {8, 4}, {10, 4}, {12, 4},
			{0, 5}, {1, 5}, {4, 5}, {5, 5}, {6, 5}, {7, 5}, {8, 5}, {11, 5}, {12, 5},
			{2, 7}, {4, 7}, {6, 7}, {8, 7}, {10, 7},
			{3, 8}, {4, 8}, {5, 8}, {7, 8}, {8, 8}, {9, 8},
		},
	},
	// Toad: A period-2 oscillator made of 6 cells.
	{
		name: "Toad",
		sparse: [][2]int{
			{1, 1}, {2, 1}, {3, 1}, {2, 2}, {3, 2}, {4, 2},
		},
	},
	// Beacon: A period-2 oscillator made of two 2x2 blocks.
	{
		name: "Beacon",
		sparse: [][2]int{
			{0, 0}, {1, 0}, {0, 1}, {1, 1}, {2, 2}, {3, 2}, {2, 3}, {3, 3},
		},
	},
	// Icon: A custom or less common pattern (please clarify if this has a standard name).
	{
		name: "Icon",
		sparse: [][2]int{
			{1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {0, 2}, {2, 2}, {4, 2}, {0, 3}, {4, 3}, {1, 4}, {2, 4}, {3, 4},
		},
	},
	// Gosper Glider Gun: The most famous pattern, periodically emits gliders indefinitely.
	{
		name: "Gosper Glider Gun",
		sparse: [][2]int{
			{1, 5}, {2, 5}, {1, 6}, {2, 6},
			{13, 3}, {14, 3}, {12, 4}, {16, 4}, {11, 5}, {17, 5}, {11, 6}, {15, 6},
			{17, 6}, {18, 6}, {11, 7}, {17, 7}, {12, 8}, {16, 8}, {13, 9}, {14, 9},
			{25, 1}, {23, 2}, {25, 2}, {21, 3}, {22, 3}, {21, 4}, {22, 4}, {21, 5},
			{22, 5}, {23, 6}, {25, 6}, {25, 7},
			{35, 3}, {36, 3}, {35, 4}, {36, 4},
		},
	},
	// Block: The simplest still life (2x2 block that never changes).
	{
		name: "Block",
		sparse: [][2]int{
			{0, 0}, {1, 0}, {0, 1}, {1, 1},
		},
	},
	// Beehive: A common still life with 6 cells.
	{
		name: "Beehive",
		sparse: [][2]int{
			{1, 0}, {2, 0}, {0, 1}, {3, 1}, {1, 2}, {2, 2},
		},
	},
	// Loaf: Another still life with 7 cells.
	{
		name: "Loaf",
		sparse: [][2]int{
			{1, 0}, {2, 0}, {0, 1}, {3, 1}, {1, 2}, {3, 2}, {2, 3},
		},
	},
	// Boat: The smallest still life with 5 cells.
	{
		name: "Boat",
		sparse: [][2]int{
			{0, 0}, {1, 0}, {0, 1}, {2, 1}, {1, 2},
		},
	},
	// Tub: A small still life with 4 cells in a diamond shape.
	{
		name: "Tub",
		sparse: [][2]int{
			{1, 0}, {0, 1}, {2, 1}, {1, 2},
		},
	},
	// Lightweight Spaceship (LWSS): A small spaceship that moves horizontally.
	{
		name: "LWSS",
		sparse: [][2]int{
			{1, 0}, {2, 0}, {3, 0}, {4, 0}, {0, 1}, {4, 1}, {4, 2}, {0, 3}, {3, 3},
		},
	},
}
