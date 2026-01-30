package size

import "fmt"

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
