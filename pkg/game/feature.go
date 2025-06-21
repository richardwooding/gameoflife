package game

type Feature struct {
	name  string
	shape [][]bool
}

func (f *Feature) GetName() string {
	return f.name
}

func (f *Feature) Stamp(cells *[][]bool, offsetX, offsetY int) {
	for y, row := range f.shape {
		for x, cell := range row {
			if cell {
				// Ensure we don't go out of bounds
				if offsetY+y < len(*cells) && offsetX+x < len((*cells)[0]) {
					(*cells)[offsetY+y][offsetX+x] = true
				}
			}
		}
	}
}

var Features = []Feature{
	{
		name: "Glider",
		shape: [][]bool{
			{false, true, false},
			{false, false, true},
			{true, true, true},
		},
	},
	{
		name: "Blinker",
		shape: [][]bool{
			{false, false, false},
			{true, true, true},
			{false, false, false},
		},
	},
	{
		name: "Pulsar",
		shape: [][]bool{
			{false, false, false, true, true, true, false, false, false},
			{false, false, true, false, false, true, false, false, false},
			{false, false, false, true, true, true, false, false, false},
			{false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false},
			{false, false, false, true, true, true, false, false, false},
			{false, false, true, false, false, true, false, false, false},
			{false, false, false, true, true, true, false, false, false},
		},
	},
	{
		name: "Toad",
		shape: [][]bool{
			{false, false, false, false, false},
			{false, true, true, true, false},
			{false, false, true, true, false},
		},
	},
	{
		name: "Beacon",
		shape: [][]bool{
			{true, true, false, false},
			{true, true, false, false},
			{false, false, true, true},
			{false, false, true, true},
		},
	},
	{
		name: "Icon",
		shape: [][]bool{
			{false, true, true, true, false},
			{true, false, false, false, true},
			{true, false, true, false, true},
			{true, false, false, false, true},
			{false, true, true, true, false},
		},
	},
	{
		name: "Gosper Glider Gun",
		shape: [][]bool{
			{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, false, true, true, false, false, false, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, true, true},
			{false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, true, true},
			{true, true, false, false, false, false, false, false, false, false, true, false, false, false, false, false, true, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
			{true, true, false, false, false, false, false, false, false, false, true, false, false, false, true, false, true, true, false, false, false, false, true, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, true, false, false, false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
			{false, false, false, false, false, false, false, false, false, false, false, false, true, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		},
	},
}
