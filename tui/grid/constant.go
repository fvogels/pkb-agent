package grid

import "pkb-agent/tui/position"

type constantGrid struct {
	cell Cell
}

func (grid *constantGrid) At(position position.Position) Cell {
	return grid.cell
}
