package game

// Feature represents a named pattern using a sparse list of live cell coordinates.
type Feature struct {
	name   string
	sparse [][2]int
}

func (f *Feature) GetName() string {
	return f.name
}

func (f *Feature) Stamp(cells *[][]bool, offsetX, offsetY int) {
	for _, pos := range f.sparse {
		x, y := pos[0]+offsetX, pos[1]+offsetY
		// Ensure we don't go out of bounds
		if y >= 0 && y < len(*cells) && x >= 0 && x < len((*cells)[0]) {
			(*cells)[y][x] = true
		}
	}
}

var Features = []Feature{
	{
		name: "Glider",
		sparse: [][2]int{
			{1, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2},
		},
	},
	{
		name: "Blinker",
		sparse: [][2]int{
			{0, 1}, {1, 1}, {2, 1},
		},
	},
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
	{
		name: "Toad",
		sparse: [][2]int{
			{1, 1}, {2, 1}, {3, 1}, {2, 2}, {3, 2}, {4, 2},
		},
	},
	{
		name: "Beacon",
		sparse: [][2]int{
			{0, 0}, {1, 0}, {0, 1}, {1, 1}, {2, 2}, {3, 2}, {2, 3}, {3, 3},
		},
	},
	{
		name: "Icon",
		sparse: [][2]int{
			{1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {0, 2}, {2, 2}, {4, 2}, {0, 3}, {4, 3}, {1, 4}, {2, 4}, {3, 4},
		},
	},
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
}
