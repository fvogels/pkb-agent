package size

import (
	"fmt"
	"pkb-agent/tui/position"
)

type Size struct {
	Width  int
	Height int
}

type Sized interface {
	Size() Size
}

func (size Size) String() string {
	return fmt.Sprintf("%dx%d", size.Width, size.Height)
}

func (size Size) IsValidPosition(position position.Position) bool {
	if position.X < 0 {
		return false
	}
	if position.Y < 0 {
		return false
	}
	if position.X >= size.Width {
		return false
	}
	if position.Y >= size.Height {
		return false
	}

	return true
}
