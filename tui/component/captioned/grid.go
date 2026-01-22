package captioned

import (
	"fmt"
	"pkb-agent/tui"
	tuigrid "pkb-agent/tui/grid"
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"
)

type grid struct {
	parent    *Component
	childGrid tuigrid.Grid
}

func newGrid(parent *Component) tuigrid.Grid {
	grid := grid{
		parent:    parent,
		childGrid: parent.child.Render(),
	}

	return &grid
}

func (grid *grid) Size() size.Size {
	childSize := grid.childGrid.Size()

	return size.Size{
		Width:  childSize.Width,
		Height: childSize.Height + 1,
	}
}

func (grid *grid) At(pos position.Position) tuigrid.Cell {
	if tui.SafeMode && !grid.isValidPosition(pos) {
		size := grid.Size()
		panic(fmt.Sprintf("invalid position %s, size %s", pos.String(), size.String()))
	}

	x := pos.X
	y := pos.Y
	caption := grid.parent.caption.Get()
	captionStyle := grid.parent.captionStyle

	if y == 0 {
		var contents rune

		if x < len(caption) {
			contents = caption[x]
		} else {
			contents = ' '
		}

		return tuigrid.Cell{
			Contents: contents,
			Style:    captionStyle,
			OnClick:  func() {},
		}
	} else {
		return grid.childGrid.At(position.Position{
			X: x,
			Y: y - 1,
		})
	}
}

func (grid *grid) isValidPosition(position position.Position) bool {
	x := position.X
	y := position.Y
	size := grid.Size()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
