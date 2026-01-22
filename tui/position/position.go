package position

import "fmt"

type Position struct {
	X int
	Y int
}

func (position Position) String() string {
	return fmt.Sprintf("(%d, %d)", position.X, position.Y)
}
