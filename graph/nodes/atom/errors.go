package atom

import (
	"fmt"
	"pkb-agent/util/pathlib"
)

type ErrMissingName struct {
	path  pathlib.Path
	index int
}

func (err *ErrMissingName) Error() string {
	return fmt.Sprintf("atom node missing name in file %s, index %d", err.path.String(), err.index)
}
