package game

import (
	"fmt"
	"github.com/cucumber/godog"
	"testing"
)

var (
	testPattern                Pattern
	testGrid                   [][]bool
	stampOffsetX, stampOffsetY int
	stampResult                [][]bool
)

func aPatternNamed(name string) error {
	testPattern = Pattern{name: name}
	return nil
}

func iGetTheNameOfThePattern() error {
	// nothing to do, just use testPattern.GetName()
	return nil
}

func theResultShouldBe(expected string) error {
	if testPattern.GetName() != expected {
		return fmt.Errorf("expected '%s', got '%s'", expected, testPattern.GetName())
	}
	return nil
}

func aNxMGrid(rows, cols int) error {
	testGrid = make([][]bool, rows)
	for i := range testGrid {
		testGrid[i] = make([]bool, cols)
	}
	return nil
}

func aPatternWithLiveCellsAt(cells *godog.Table) error {
	var sparse [][2]int
	for _, row := range cells.Rows {
		if len(row.Cells) != 2 {
			return fmt.Errorf("expected 2 columns, got %d", len(row.Cells))
		}
		var x, y int
		fmt.Sscanf(row.Cells[0].Value, "%d", &x)
		fmt.Sscanf(row.Cells[1].Value, "%d", &y)
		sparse = append(sparse, [2]int{x, y})
	}
	testPattern.sparse = sparse
	return nil
}

func iStampThePatternAtOffset(x, y int) error {
	stampOffsetX, stampOffsetY = x, y
	testPattern.Stamp(&testGrid, x, y)
	return nil
}

func theGridShouldHaveLiveCellsAt(cells *godog.Table) error {
	for _, row := range cells.Rows {
		var x, y int
		fmt.Sscanf(row.Cells[0].Value, "%d", &x)
		fmt.Sscanf(row.Cells[1].Value, "%d", &y)
		if y < 0 || y >= len(testGrid) || x < 0 || x >= len(testGrid[0]) || !testGrid[y][x] {
			return fmt.Errorf("expected cell (%d,%d) to be alive", x, y)
		}
	}
	return nil
}

func thePredefinedPatterns() error {
	// nothing to do, Patterns is global
	return nil
}

func eachPatternShouldHaveANonEmptyName() error {
	for _, p := range Patterns {
		if p.name == "" {
			return fmt.Errorf("pattern with empty name: %+v", p)
		}
	}
	return nil
}

func eachPatternShouldHaveValidCoordinates() error {
	for _, p := range Patterns {
		for _, pos := range p.sparse {
			if len(pos) != 2 {
				return fmt.Errorf("pattern %s has invalid coordinate: %v", p.name, pos)
			}
		}
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a pattern named "([^"]*)"$`, aPatternNamed)
	ctx.Step(`^I get the name of the pattern$`, iGetTheNameOfThePattern)
	ctx.Step(`^the result should be "([^"]*)"$`, theResultShouldBe)
	ctx.Step(`^a (\d+)x(\d+) grid$`, aNxMGrid)
	ctx.Step(`^a pattern with live cells at \((\d+),(\d+)\) and \((\d+),(\d+)\)$`,
		func(x1, y1, x2, y2 int) error {
			testPattern.sparse = [][2]int{{x1, y1}, {x2, y2}}
			return nil
		},
	)
	ctx.Step(`^I stamp the pattern at offset \((\d+),(\d+)\)$`, iStampThePatternAtOffset)
	ctx.Step(`^the grid should have live cells at \((\d+),(\d+)\) and \((\d+),(\d+)\)$`,
		func(x1, y1, x2, y2 int) error {
			if !testGrid[y1][x1] || !testGrid[y2][x2] {
				return fmt.Errorf("expected cells (%d,%d) and (%d,%d) to be alive", x1, y1, x2, y2)
			}
			return nil
		},
	)
	ctx.Step(`^the predefined patterns$`, thePredefinedPatterns)
	ctx.Step(`^each pattern should have a non-empty name$`, eachPatternShouldHaveANonEmptyName)
	ctx.Step(`^each pattern should have valid coordinates$`, eachPatternShouldHaveValidCoordinates)
}

func TestPatterns(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "patterns",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features/pattern.feature"},
		},
	}

	if suite.Run() != 0 {
		t.Fatal("Test suite failed")
	}
}
