package stringsview

import (
	"pkb-agent/tui"
)

type grid struct {
	parent *Component
}

func newGrid(component *Component) *grid {
	return &grid{
		parent: component,
	}
}

func (grid *grid) GetSize() tui.Size {
	return grid.parent.size
}

func (grid *grid) Get(position tui.Position) tui.Cell {
	if tui.SafeMode && !grid.isValidPosition(position) {
		panic("invalid coordinates")
	}

	x := position.X
	y := position.Y
	firstVisibleIndex := grid.parent.firstVisibleIndex
	currentItemIndex := firstVisibleIndex + y
	items := grid.parent.items

	var contents rune
	var style *tui.Style

	if currentItemIndex >= items.Size() {
		// Current line is outside of bounds of list
		contents = ' '
		style = grid.parent.emptyStyle
	} else {
		// Current line contains item
		currentItem := items.At(currentItemIndex)

		if x < len(currentItem.Runes) {
			contents = currentItem.Runes[x]
		} else {
			contents = ' '
		}

		style = currentItem.Style
	}

	cell := tui.Cell{
		Contents: contents,
		Style:    style,
		OnClick:  nil,
	}

	return cell
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.GetSize()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
