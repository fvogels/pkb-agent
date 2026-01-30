package grid

import (
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"
)

type cropper struct {
	grid Grid
	size size.Size
}

func (grid *cropper) Size() size.Size {
	return grid.size
}

func (grid *cropper) At(position position.Position) Cell {
	if !grid.size.IsValidPosition(position) {
		panic("invalid position")
	}

	return grid.At(position)
}
