package ansiview

import (
	"pkb-agent/tui"
	tuigrid "pkb-agent/tui/grid"
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"
)

type grid struct {
	size       size.Size
	ansiGrid   tuigrid.FiniteGrid
	emptyStyle *tui.Style
}

func newGrid(size size.Size, ansiGrid tuigrid.FiniteGrid, emptyStyle *tui.Style) tuigrid.FiniteGrid {
	return &grid{
		size:       size,
		ansiGrid:   ansiGrid,
		emptyStyle: emptyStyle,
	}
}

func (graph *grid) Size() size.Size {
	return graph.size
}

func (graph *grid) At(position position.Position) tuigrid.Cell {
	x := position.X
	y := position.Y
	gridSize := graph.ansiGrid.Size()

	if x < gridSize.Width && y < gridSize.Height {
		return graph.ansiGrid.At(position)
	} else {
		return tuigrid.Cell{
			Contents: ' ',
			Style:    graph.emptyStyle,
		}
	}
}
