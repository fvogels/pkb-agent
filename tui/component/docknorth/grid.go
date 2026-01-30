package docknorth

import (
	tuigrid "pkb-agent/tui/grid"
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"
)

type grid struct {
	size            size.Size
	mainChildGrid   tuigrid.FiniteGrid
	dockedChildGrid tuigrid.FiniteGrid
	boundary        int // Y-coordinate of where docked child ends
}

func (grid *grid) Size() size.Size {
	return grid.size
}

func (grid *grid) At(pos position.Position) tuigrid.Cell {
	if !grid.isValidPosition(pos) {
		panic("invalid coordinates")
	}

	if pos.Y < grid.boundary {
		return grid.dockedChildGrid.At(pos)
	} else {
		return grid.mainChildGrid.At(position.Position{
			X: pos.X,
			Y: pos.Y - grid.boundary,
		})
	}
}

func (grid *grid) isValidPosition(position position.Position) bool {
	return grid.Size().IsValidPosition(position)
}
