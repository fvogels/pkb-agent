package stringlist

import (
	"pkb-agent/tui"
)

type grid struct {
	size          tui.Size
	items         [][]rune
	selectedIndex int
	emptyStyle    *tui.Style
	itemStyle     *tui.Style
	selectedStyle *tui.Style
}

func (grid *grid) GetSize() tui.Size {
	return grid.size
}

func (grid *grid) Get(position tui.Position) tui.Cell {
	if tui.SafeMode && !grid.isValidPosition(position) {
		panic("invalid coordinates")
	}

	x := position.X
	y := position.Y
	selectedIndex := grid.selectedIndex
	items := grid.items

	var contents rune
	var style *tui.Style

	if y >= len(items) {
		// Current line is outside of bounds of list
		contents = ' '
		style = grid.emptyStyle
	} else {
		// Current line contains item
		visibleItem := items[y]

		if x < len(visibleItem) {
			contents = visibleItem[x]
		} else {
			contents = ' '
		}

		if y == selectedIndex {
			style = grid.selectedStyle
		} else {
			style = grid.itemStyle
		}
	}

	cell := tui.Cell{
		Contents: contents,
		Style:    style,
	}

	return cell
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.GetSize()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
