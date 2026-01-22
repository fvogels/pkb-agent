package ansiview

import (
	"pkb-agent/tui"
	"pkb-agent/tui/position"
)

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

func (graph *grid) Size() tui.Size {
	return graph.size
}

func (graph *grid) At(position position.Position) tui.Cell {
	x := position.X
	y := position.Y
	gridSize := graph.ansiGrid.Size()

	if x < gridSize.Width && y < gridSize.Height {
		return graph.ansiGrid.At(position)
	} else {
		return tui.Cell{
			Contents: ' ',
			Style:    graph.emptyStyle,
		}
	}
}
