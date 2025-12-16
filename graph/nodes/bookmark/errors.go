package bookmark

import (
	"fmt"
	"pkb-agent/util/pathlib"
)

type ErrMissingName struct {
	path  pathlib.Path
	index int
}

type ErrMissingDescription struct {
	path  pathlib.Path
	index int
}

func (err *ErrMissingName) Error() string {
	return fmt.Sprintf("bookmark node missing name in file %s, index %d", err.path.String(), err.index)
}

func (err *ErrMissingDescription) Error() string {
	return fmt.Sprintf("bookmark node missing description in file %s, index %d", err.path.String(), err.index)
}
