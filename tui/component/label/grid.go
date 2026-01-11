package label

import (
	"fmt"
	"pkb-agent/tui"
)

type grid struct {
	parent   *Component
	contents []rune
}

func newGrid(parent *Component, contents []rune) *grid {
	result := grid{
		parent:   parent,
		contents: contents,
	}

	return &result
}

func (grid *grid) GetSize() tui.Size {
	return grid.parent.size
}

func (grid *grid) Get(position tui.Position) tui.Cell {
	if tui.SafeMode && !grid.isValidPosition(position) {
		size := grid.parent.size
		panic(fmt.Sprintf("invalid grid access: parent: %s, coordinates (%d, %d), size: %dx%d, contents: %s", grid.parent.name, position.X, position.Y, size.Width, size.Height, string(grid.contents)))
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
		Style:    grid.parent.style,
	}

	return cell
}

func (grid *grid) isValidPosition(position tui.Position) bool {
	x := position.X
	y := position.Y
	size := grid.GetSize()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
