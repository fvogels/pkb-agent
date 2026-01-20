package docksouth

import "pkb-agent/tui"

type grid struct {
	size            tui.Size
	mainChildGrid   tui.Grid
	dockedChildGrid tui.Grid
	boundary        int // Y-coordinate of where docked child starts
}

func (grid *grid) Size() tui.Size {
	return grid.size
}

func (grid *grid) At(position tui.Position) tui.Cell {
	if tui.SafeMode && !grid.isValidPosition(position) {
		panic("invalid coordinates")
	}

	if position.Y < grid.boundary {
		return grid.mainChildGrid.At(position)
	} else {
		return grid.dockedChildGrid.At(tui.Position{
			X: position.X,
			Y: position.Y - grid.boundary,
		})
	}
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.Size()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
