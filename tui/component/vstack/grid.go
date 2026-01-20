package vstack

import (
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
)

type grid struct {
	parent     *Component
	childGrids list.List[tui.Grid]
}

func newGrid(parent *Component, childGrids list.List[tui.Grid]) tui.Grid {
	return &grid{
		parent:     parent,
		childGrids: childGrids,
	}
}

func (grid *grid) GetSize() tui.Size {
	return grid.parent.Size
}

func (grid *grid) Get(position tui.Position) tui.Cell {
	if !grid.isValidPosition(position) {
		panic("invalid position")
	}

	x := position.X
	y := position.Y
	i := 0

	for i < grid.childGrids.Size() && y >= grid.childGrids.At(i).GetSize().Height {
		y -= grid.childGrids.At(i).GetSize().Height
		i++
	}

	if i == grid.childGrids.Size() {
		return tui.Cell{
			Contents: ' ',
			Style:    grid.parent.emptyStyle,
			OnClick:  nil,
		}
	} else {
		return grid.childGrids.At(i).Get(tui.Position{X: x, Y: y})
	}
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.parent.Size
	width := size.Width
	height := size.Height

	return 0 <= x && x < width && 0 <= y && y < height
}
