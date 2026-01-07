package strlist

import (
	"pkb-agent/tui"
)

type grid struct {
	size  tui.Size
	cells []tui.Cell
}

func NewGrid(parent *Component) *grid {
	width := parent.size.Width
	height := parent.size.Height
	items := parent.items
	cells := make([]tui.Cell, width*height)

	for y := range height {
		row := cells[y*width : (y+1)*width]
		itemIndex := parent.firstVisibleIndex + y

		if itemIndex < items.Size() {
			item := []rune(items.At(itemIndex))

			var style *tui.Style
			if itemIndex == parent.selectedIndex.Get() {
				style = parent.selectedStyle
			} else {
				style = parent.itemStyle
			}

			x := 0
			for x < len(item) && x < width {
				row[x].Contents = item[x]
				row[x].Style = style
				x++
			}

			for x < width {
				row[x].Contents = ' '
				row[x].Style = style
				x++
			}
		} else {
			// Beyond end of items; fill in empty line
			for x := range width {
				row[x].Contents = ' '
				row[x].Style = parent.emptyStyle
			}
		}
	}

	return &grid{
		size:  parent.size,
		cells: cells,
	}
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
	width := grid.size.Width

	return grid.cells[y*width+x]
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.GetSize()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
