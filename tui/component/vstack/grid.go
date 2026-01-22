package vstack

import (
	"pkb-agent/persistent/list"
	tuigrid "pkb-agent/tui/grid"
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"
)

type grid struct {
	parent     *Component
	childGrids list.List[tuigrid.Grid]
}

func newGrid(parent *Component, childGrids list.List[tuigrid.Grid]) tuigrid.Grid {
	return &grid{
		parent:     parent,
		childGrids: childGrids,
	}
}

func (grid *grid) Size() size.Size {
	return grid.parent.Size
}

func (grid *grid) At(pos position.Position) tuigrid.Cell {
	if !grid.isValidPosition(pos) {
		panic("invalid position")
	}

	x := pos.X
	y := pos.Y
	i := 0

	for i < grid.childGrids.Size() && y >= grid.childGrids.At(i).Size().Height {
		y -= grid.childGrids.At(i).Size().Height
		i++
	}

	if i == grid.childGrids.Size() {
		return tuigrid.Cell{
			Contents: ' ',
			Style:    grid.parent.emptyStyle,
			OnClick:  nil,
		}
	} else {
		return grid.childGrids.At(i).At(position.Position{X: x, Y: y})
	}
}

func (grid *grid) isValidPosition(pos position.Position) bool {
	x := pos.X
	y := pos.Y
	size := grid.parent.Size
	width := size.Width
	height := size.Height

	return 0 <= x && x < width && 0 <= y && y < height
}
