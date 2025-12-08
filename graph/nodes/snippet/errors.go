package snippet

import (
	"fmt"
	"pkb-agent/util/pathlib"
)

type ErrMissingSnippet struct {
	path pathlib.Path
}

func (err *ErrMissingSnippet) Error() string {
	return fmt.Sprintf("only metadata in snippet in %s", err.path.String())
}

type ErrMissingName struct {
	path pathlib.Path
}

func (err *ErrMissingName) Error() string {
	return fmt.Sprintf("node in %s is missing name", err.path.String())
}
