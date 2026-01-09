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

func (graph *grid) GetSize() tui.Size {
	return graph.size
}

func (graph *grid) Get(position tui.Position) tui.Cell {
	x := position.X
	y := position.Y
	gridSize := graph.ansiGrid.GetSize()

	if x < gridSize.Width && y < gridSize.Height {
		return graph.ansiGrid.Get(position)
	} else {
		return tui.Cell{
			Contents: ' ',
			Style:    graph.emptyStyle,
		}
	}
}
