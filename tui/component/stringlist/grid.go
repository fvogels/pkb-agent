package stringlist

import (
	"pkb-agent/tui"
)

type grid struct {
	parent           *Component
	rows             [][]rune
	rowClickHandlers []func()
	selectedRowIndex int
}

func newGrid(component *Component) *grid {
	lineCount := component.size.Height
	itemCount := component.items.Size()
	itemIndex := component.firstVisibleIndex
	lineIndex := 0
	itemsAsRunes := make([][]rune, lineCount)
	rowClickHandlers := make([]func(), lineCount)

	for lineIndex < lineCount && itemIndex < itemCount {
		item := component.items.At(itemIndex)
		itemsAsRunes[lineIndex] = []rune(item)

		itemIndexCopy := itemIndex
		rowClickHandlers[lineIndex] = func() {
			component.onSelectionChanged(itemIndexCopy)
		}
		lineIndex++
		itemIndex++
	}

	return &grid{
		parent:           component,
		rows:             itemsAsRunes,
		selectedRowIndex: component.selectedIndex.Get() - component.firstVisibleIndex,
		rowClickHandlers: rowClickHandlers,
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
	selectedRowIndex := grid.selectedRowIndex
	items := grid.rows

	var contents rune
	var style *tui.Style

	if y >= len(items) {
		// Current line is outside of bounds of list
		contents = ' '
		style = grid.parent.emptyStyle
	} else {
		// Current line contains item
		visibleItem := items[y]

		if x < len(visibleItem) {
			contents = visibleItem[x]
		} else {
			contents = ' '
		}

		if y == selectedRowIndex {
			style = grid.parent.selectedItemStyle
		} else {
			style = grid.parent.itemStyle
		}
	}

	cell := tui.Cell{
		Contents: contents,
		Style:    style,
		OnClick:  grid.rowClickHandlers[y],
	}

	return cell
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.GetSize()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
