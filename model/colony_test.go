package model

import (
	"fmt"
	"github.com/cucumber/godog"
	"testing"
)

type colonyFeature struct {
	colony *Colony
}

func (f *colonyFeature) aColonyOfSize(arg1, arg2 int) error {
	c := make([][]bool, arg2)
	for i := range c {
		c[i] = make([]bool, arg1)
	}
	f.colony = &Colony{}
	f.colony.SetCells(c)
	f.colony.dx = arg1
	f.colony.dy = arg2
	return nil
}

func (f *colonyFeature) theCellAtIsAlive(x, y int) error {
	(*f.colony.cells)[y][x] = true
	return nil
}

func (f *colonyFeature) nextGenerationIsComputed() error {
	f.colony.Generate()
	return nil
}

func (f *colonyFeature) theCellAtShouldBeState(x, y int, state string) error {
	alive := (*f.colony.cells)[y][x]
	shouldBeAlive := state == "alive"
	if alive != shouldBeAlive {
		return fmt.Errorf("expected cell (%d,%d) to be %s, but it was %v", x, y, state, alive)
	}
	return nil
}

func (f *colonyFeature) toggleCellAt(x, y int) error {
	if x < 0 || y < 0 || x >= f.colony.dx || y >= f.colony.dy {
		return fmt.Errorf("cell (%d,%d) is out of bounds", x, y)
	}
	(*f.colony.cells)[y][x] = !(*f.colony.cells)[y][x]
	return nil
}

func (f *colonyFeature) allCellsShouldBeDead() error {
	for y := 0; y < f.colony.dy; y++ {
		for x := 0; x < f.colony.dx; x++ {
			if (*f.colony.cells)[y][x] {
				return fmt.Errorf("cell (%d,%d) is alive, expected all cells to be dead", x, y)
			}
		}
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	f := &colonyFeature{}
	ctx.Step(`^a (\d+)x(\d+) colony$`, f.aColonyOfSize)
	ctx.Step(`^the cell at \((\d+),(\d+)\) is alive$`, f.theCellAtIsAlive)
	ctx.Step(`^the next generation is computed$`, f.nextGenerationIsComputed)
	ctx.Step(`^the cell at \((\d+),(\d+)\) should be (alive|dead)$`, f.theCellAtShouldBeState)
	ctx.Step(`^I toggle the cell at \((\d+),(\d+)\)$`, f.toggleCellAt)
	ctx.Step(`^all cells should be dead$`, f.allCellsShouldBeDead)
}

func TestColonyFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "colony",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features/colony.feature"},
		},
	}

	if suite.Run() != 0 {
		t.Fatal("there were test failures")
	}
}
