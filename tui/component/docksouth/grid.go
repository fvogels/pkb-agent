package docksouth

import "pkb-agent/tui"

type grid struct {
	size            tui.Size
	mainChildGrid   tui.Grid
	dockedChildGrid tui.Grid
	boundary        int // Y-coordinate of where docked child starts
}

func (grid *grid) GetSize() tui.Size {
	return grid.size
}

func (grid *grid) Get(position tui.Position) tui.Cell {
	if tui.SafeMode && !grid.isValidPosition(position) {
		panic("invalid coordinates")
	}

	if position.Y < grid.boundary {
		return grid.mainChildGrid.Get(position)
	} else {
		return grid.dockedChildGrid.Get(tui.Position{
			X: position.X,
			Y: position.Y - grid.boundary,
		})
	}
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.GetSize()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
