package grid

import (
	"fmt"
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"
)

type materializedGrid struct {
	items []Cell
	size  size.Size
}

func Materialize(g FiniteGrid) FiniteGrid {
	size := g.Size()
	items := make([]Cell, size.Width*size.Height)

	i := 0
	for y := range size.Height {
		for x := range size.Width {
			position := position.Position{X: x, Y: y}
			items[i] = g.At(position)

			i++
		}
	}

	return &materializedGrid{
		items: items,
		size:  size,
	}
}

func NewMaterializedGrid(size size.Size, initializer func(position.Position) Cell) *materializedGrid {
	items := make([]Cell, size.Width*size.Height)

	i := 0
	for y := range size.Height {
		for x := range size.Width {
			position := position.Position{X: x, Y: y}
			items[i] = initializer(position)
			i++
		}
	}

	result := materializedGrid{
		items: items,
		size:  size,
	}

	return &result
}

func (grid *materializedGrid) Size() size.Size {
	return grid.size
}

func (grid *materializedGrid) At(position position.Position) Cell {
	if !grid.isValidPosition(position) {
		panic(fmt.Sprintf("invalid position (%d, %d), size %dx%d", position.X, position.Y, grid.size.Width, grid.size.Height))
	}

	return grid.items[grid.computeIndexOfPosition(position)]
}

func (grid *materializedGrid) Set(position position.Position, cell Cell) {
	if !grid.isValidPosition(position) {
		panic("invalid position")
	}

	grid.items[grid.computeIndexOfPosition(position)] = cell
}

func (grid *materializedGrid) isValidPosition(position position.Position) bool {
	return grid.Size().IsValidPosition(position)
}

func (grid *materializedGrid) computeIndexOfPosition(position position.Position) int {
	width := grid.size.Width
	x := position.X
	y := position.Y

	return y*width + x
}
