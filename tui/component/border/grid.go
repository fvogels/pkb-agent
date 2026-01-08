package border

import (
	"pkb-agent/tui"

	"github.com/gdamore/tcell/v3"
)

type grid struct {
	childGrid tui.Grid
	style     tui.Style
}

func (grid *grid) GetSize() tui.Size {
	childSize := grid.childGrid.GetSize()

	return tui.Size{
		Width:  childSize.Width + 2,
		Height: childSize.Height + 2,
	}
}

func (grid *grid) Get(position tui.Position) tui.Cell {
	if tui.SafeMode && !grid.isValidPosition(position) {
		panic("invalid coordinates")
	}

	x := position.X
	y := position.Y
	size := grid.GetSize()
	width := size.Width
	height := size.Height

	var char rune
	var style *tui.Style
	var onClick func()

	if x == 0 {
		style = &grid.style

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
		style = &grid.style

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
		style = &grid.style
		char = tcell.RuneHLine
	} else {
		cell := grid.childGrid.Get(tui.Position{X: x - 1, Y: y - 1})
		style = cell.Style
		char = cell.Contents
		onClick = cell.OnClick
	}

	return tui.Cell{
		Contents: char,
		Style:    style,
		OnClick:  onClick,
	}
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.GetSize()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
