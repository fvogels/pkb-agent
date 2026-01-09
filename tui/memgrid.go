package tui

type MemoryGrid struct {
	items []Cell
	size  Size
}

func MaterializeGrid(grid Grid) Grid {
	size := grid.GetSize()
	items := make([]Cell, size.Width*size.Height)

	i := 0
	for y := range size.Height {
		for x := range size.Width {
			position := Position{X: x, Y: y}
			items[i] = grid.Get(position)

			i++
		}
	}

	return &MemoryGrid{
		items: items,
		size:  size,
	}
}

func (grid *MemoryGrid) GetSize() Size {
	return grid.size
}

func (grid *MemoryGrid) Get(position Position) Cell {
	if SafeMode && !grid.isValidPosition(position) {
		panic("invalid position")
	}

	width := grid.size.Width
	x := position.X
	y := position.Y

	return grid.items[y*width+x]
}

func (grid *MemoryGrid) isValidPosition(position Position) bool {
	if position.X < 0 {
		return false
	}
	if position.Y < 0 {
		return false
	}
	if position.X >= grid.size.Width {
		return false
	}
	if position.Y >= grid.size.Height {
		return false
	}

	return true
}
