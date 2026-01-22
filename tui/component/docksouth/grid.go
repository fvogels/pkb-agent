package docksouth

import (
	"pkb-agent/tui"
	"pkb-agent/tui/position"
)

type grid struct {
	size            tui.Size
	mainChildGrid   tui.Grid
	dockedChildGrid tui.Grid
	boundary        int // Y-coordinate of where docked child starts
}

func (grid *grid) Size() tui.Size {
	return grid.size
}

func (grid *grid) At(pos position.Position) tui.Cell {
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
