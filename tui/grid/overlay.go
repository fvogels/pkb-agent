package grid

import (
	"pkb-agent/tui/position"
)

type overlay struct {
	layers []FiniteGrid
}

func OverlayGrids(layers ...FiniteGrid) Grid {
	return &overlay{
		layers: layers,
	}
}

func (grid *overlay) At(pos position.Position) Cell {
	for _, layer := range grid.layers {
		cell := layer.At(pos)

		if cell.Contents != 0 {
			return cell
		}
	}

	return Cell{}
}
