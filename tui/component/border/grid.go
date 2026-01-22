package border

import (
	"fmt"
	"pkb-agent/tui"
	tuigrid "pkb-agent/tui/grid"
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"

	"github.com/gdamore/tcell/v3"
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
		Width:  childSize.Width + 2,
		Height: childSize.Height + 2,
	}
}

func (g *grid) At(pos position.Position) tuigrid.Cell {
	if tui.SafeMode && !g.isValidPosition(pos) {
		size := g.Size()
		panic(fmt.Sprintf("invalid position (%d, %d), size %dx%d in component %s", pos.X, pos.Y, size.Width, size.Height, g.parent.Name))
	}

	x := pos.X
	y := pos.Y
	size := g.Size()
	width := size.Width
	height := size.Height

	var char rune
	var style *tui.Style
	var onClick func()

	if x == 0 {
		style = g.parent.style

		if y == 0 {
			// Upper left corner
			char = tcell.RuneULCorner
		} else if y == height-1 {
			// Lower left corner
			char = tcell.RuneLLCorner
		} else {
			// Left border
			char = tcell.RuneVLine
		}
	} else if x == width-1 {
		style = g.parent.style

		if y == 0 {
			// Upper right corner
			char = tcell.RuneURCorner
		} else if y == height-1 {
			// Lower right corner
			char = tcell.RuneLRCorner
		} else {
			// Right border
			char = tcell.RuneVLine
		}
	} else if y == 0 || y == height-1 {
		style = g.parent.style
		char = tcell.RuneHLine
	} else {
		cell := g.childGrid.At(position.Position{X: x - 1, Y: y - 1})
		style = cell.Style
		char = cell.Contents
		onClick = cell.OnClick
	}

	return tuigrid.Cell{
		Contents: char,
		Style:    style,
		OnClick:  onClick,
	}
}

func (g *grid) isValidPosition(position position.Position) bool {
	x := position.X
	y := position.Y
	size := g.Size()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
