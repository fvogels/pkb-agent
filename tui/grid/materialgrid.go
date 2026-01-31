package grid

import (
	"fmt"
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"
)

type MaterializedGrid struct {
	items []Cell
	size  size.Size
}

func Materialize(g FiniteGrid) *MaterializedGrid {
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

	return &MaterializedGrid{
		items: items,
		size:  size,
	}
}

func NewMaterializedGrid(size size.Size, initializer func(position.Position) Cell) *MaterializedGrid {
	items := make([]Cell, size.Width*size.Height)

	i := 0
	for y := range size.Height {
		for x := range size.Width {
			position := position.Position{X: x, Y: y}
			items[i] = initializer(position)
			i++
		}
	}

	result := MaterializedGrid{
		items: items,
		size:  size,
	}

	return &result
}

func (grid *MaterializedGrid) Size() size.Size {
	return grid.size
}

func (grid *MaterializedGrid) At(position position.Position) Cell {
	if !grid.isValidPosition(position) {
		panic(fmt.Sprintf("invalid position (%d, %d), size %dx%d", position.X, position.Y, grid.size.Width, grid.size.Height))
	}

	return grid.items[grid.computeIndexOfPosition(position)]
}

func (grid *MaterializedGrid) Set(position position.Position, cell Cell) {
	if !grid.isValidPosition(position) {
		panic("invalid position")
	}

	grid.items[grid.computeIndexOfPosition(position)] = cell
}

func (grid *MaterializedGrid) isValidPosition(position position.Position) bool {
	return grid.Size().IsValidPosition(position)
}

func (grid *MaterializedGrid) computeIndexOfPosition(position position.Position) int {
	width := grid.size.Width
	x := position.X
	y := position.Y

	return y*width + x
}
