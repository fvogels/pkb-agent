package graph

import "fmt"

type ErrNameClash struct {
	name string
}

func (err *ErrNameClash) Error() string {
	return fmt.Sprintf("name clash: %s", err.name)
}

type ErrUnknownNodes struct{}

func (err *ErrUnknownNodes) Error() string {
	return "unknown node"
}
