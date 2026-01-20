package captioned

import (
	"fmt"
	"pkb-agent/tui"
)

type grid struct {
	parent    *Component
	childGrid tui.Grid
}

func newGrid(parent *Component) tui.Grid {
	grid := grid{
		parent:    parent,
		childGrid: parent.child.Render(),
	}

	return &grid
}

func (grid *grid) GetSize() tui.Size {
	childSize := grid.childGrid.GetSize()

	return tui.Size{
		Width:  childSize.Width,
		Height: childSize.Height + 1,
	}
}

func (grid *grid) Get(position tui.Position) tui.Cell {
	if tui.SafeMode && !grid.isValidPosition(position) {
		size := grid.GetSize()
		panic(fmt.Sprintf("invalid position %s, size %s", position.String(), size.String()))
	}

	x := position.X
	y := position.Y
	caption := grid.parent.caption.Get()
	captionStyle := grid.parent.captionStyle

	if y == 0 {
		var contents rune

		if x < len(caption) {
			contents = caption[x]
		} else {
			contents = ' '
		}

		return tui.Cell{
			Contents: contents,
			Style:    captionStyle,
			OnClick:  func() {},
		}
	} else {
		return grid.childGrid.Get(tui.Position{
			X: x,
			Y: y - 1,
		})
	}
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.GetSize()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
