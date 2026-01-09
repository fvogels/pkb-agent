package ansiview

import "pkb-agent/tui"

type grid struct {
	size       tui.Size
	ansiGrid   tui.Grid
	emptyStyle *tui.Style
}

func newGrid(size tui.Size, ansiGrid tui.Grid, emptyStyle *tui.Style) tui.Grid {
	return &grid{
		size:       size,
		ansiGrid:   ansiGrid,
		emptyStyle: emptyStyle,
	}
}

func (g *grid) GetSize() tui.Size {
	return g.size
}

func (grid *grid) Get(position tui.Position) tui.Cell {
	x := position.X
	y := position.Y
	gridSize := grid.ansiGrid.GetSize()

	if x < gridSize.Width && y < gridSize.Height {
		return grid.ansiGrid.Get(position)
	} else {
		return tui.Cell{
			Contents: ' ',
			Style:    grid.emptyStyle,
		}
	}
}
