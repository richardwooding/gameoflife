package main

import (
	"github.com/richardwooding/gameoflife/pkg/life"
	"os"
)

func main() {
	life := life.New(5, 5)
	life.Alive(2, 1)
	life.Alive(2, 2)
	life.Alive(2, 3)
	life.Print(os.Stdout)
	life.GenerateAndPrint(os.Stdout)
	life.GenerateAndPrint(os.Stdout)
	life.GenerateAndPrint(os.Stdout)
	life.GenerateAndPrint(os.Stdout)
}
