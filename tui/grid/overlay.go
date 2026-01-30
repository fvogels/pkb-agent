package grid

import (
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"
)

type overlay struct {
	size   size.Size
	layers []FiniteGrid
}

func OverlayGrids(size size.Size, layers ...FiniteGrid) FiniteGrid {
	return &overlay{
		size:   size,
		layers: layers,
	}
}

func (grid *overlay) Size() size.Size {
	return grid.size
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

func (grid *overlay) isValidPosition(position position.Position) bool {
	return grid.Size().IsValidPosition(position)
}
