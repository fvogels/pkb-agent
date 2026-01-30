package grid

import "pkb-agent/tui/position"

type positionedGrid struct {
	grid Grid
	dx   int
	dy   int
}

func Reposition(grid Grid, dx int, dy int) Grid {
	return &positionedGrid{
		grid: grid,
		dx:   dx,
		dy:   dy,
	}
}

func (grid *positionedGrid) At(position position.Position) Cell {
	return grid.grid.At(position.Move(grid.dx, grid.dy))
}
