package hybrid

import (
	"fmt"
	"pkb-agent/util/pathlib"
)

type ErrMissingName struct {
	path  pathlib.Path
	index int
}

type ErrInvalidAction struct {
	path pathlib.Path
}

func (err *ErrMissingName) Error() string {
	return fmt.Sprintf("bookmark node missing name in file %s, index %d", err.path.String(), err.index)
}

func (err *ErrInvalidAction) Error() string {
	return "invalid node action"
}
