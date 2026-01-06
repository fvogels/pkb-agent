package label

import (
	"pkb-agent/tui"
)

type grid struct {
	contents []rune
	style    *tui.Style
	size     tui.Size
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

	var contents rune
	if x < len(grid.contents) && y == 0 {
		contents = grid.contents[x]
	} else {
		contents = ' '
	}

	cell := tui.Cell{
		Contents: contents,
		Style:    grid.style,
	}

	return cell
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.GetSize()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
