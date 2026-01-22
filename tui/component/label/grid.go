package label

import (
	"fmt"
	"pkb-agent/tui"
	"pkb-agent/tui/position"
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

func (grid *grid) Size() tui.Size {
	return grid.parent.Size
}

func (grid *grid) At(position position.Position) tui.Cell {
	if tui.SafeMode && !grid.isValidPosition(position) {
		size := grid.parent.Size
		panic(fmt.Sprintf("invalid grid access: parent: %s, coordinates (%d, %d), size: %dx%d, contents: %s", grid.parent.Name, position.X, position.Y, size.Width, size.Height, string(grid.contents)))
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

func (grid *grid) isValidPosition(position position.Position) bool {
	x := position.X
	y := position.Y
	size := grid.Size()

	return 0 <= x && x < size.Width && 0 <= y && y < size.Height
}
