package tui

import "fmt"

type MemoryGrid struct {
	items []Cell
	size  Size
}

func MaterializeGrid(grid Grid) Grid {
	size := grid.Size()
	items := make([]Cell, size.Width*size.Height)

	i := 0
	for y := range size.Height {
		for x := range size.Width {
			position := Position{X: x, Y: y}
			items[i] = grid.At(position)

			i++
		}
	}

	return &MemoryGrid{
		items: items,
		size:  size,
	}
}

func NewMaterializedGrid(size Size, initializer func(Position) Cell) *MemoryGrid {
	items := make([]Cell, size.Width*size.Height)

	i := 0
	for y := range size.Height {
		for x := range size.Width {
			position := Position{X: x, Y: y}
			items[i] = initializer(position)
			i++
		}
	}

	result := MemoryGrid{
		items: items,
		size:  size,
	}

	return &result
}

func (grid *MemoryGrid) Size() Size {
	return grid.size
}

func (grid *MemoryGrid) At(position Position) Cell {
	if SafeMode && !grid.isValidPosition(position) {
		panic(fmt.Sprintf("invalid position (%d, %d), size %dx%d", position.X, position.Y, grid.size.Width, grid.size.Height))
	}

	return grid.items[grid.computeIndexOfPosition(position)]
}

func (grid *MemoryGrid) Set(position Position, cell Cell) {
	if SafeMode && !grid.isValidPosition(position) {
		panic("invalid position")
	}

	grid.items[grid.computeIndexOfPosition(position)] = cell
}

func (grid *MemoryGrid) isValidPosition(position Position) bool {
	if position.X < 0 {
		return false
	}
	if position.Y < 0 {
		return false
	}
	if position.X >= grid.size.Width {
		return false
	}
	if position.Y >= grid.size.Height {
		return false
	}

	return true
}

func (grid *MemoryGrid) computeIndexOfPosition(position Position) int {
	width := grid.size.Width
	x := position.X
	y := position.Y

	return y*width + x
}
