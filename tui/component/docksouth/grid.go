package docksouth

import (
	"pkb-agent/tui"
	tuigrid "pkb-agent/tui/grid"
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"
)

type grid struct {
	size            size.Size
	mainChildGrid   tuigrid.Grid
	dockedChildGrid tuigrid.Grid
	boundary        int // Y-coordinate of where docked child starts
}

func (grid *grid) Size() size.Size {
	return grid.size
}

func (grid *grid) At(pos position.Position) tuigrid.Cell {
	if tui.SafeMode && !grid.isValidPosition(pos) {
		panic("invalid coordinates")
	}

	if pos.Y < grid.boundary {
		return grid.mainChildGrid.At(pos)
	} else {
		return grid.dockedChildGrid.At(position.Position{
			X: pos.X,
			Y: pos.Y - grid.boundary,
		})
	}
}

func (grid *grid) isValidPosition(pos position.Position) bool {
	x := pos.X
	y := pos.Y
	size := grid.Size()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
