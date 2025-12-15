package markdown

import (
	"fmt"
	"pkb-agent/util/pathlib"
)

type ErrMissingSnippet struct {
	path pathlib.Path
}

type ErrMissingName struct{}

type ErrMalformed struct {
	path pathlib.Path
}

func (err *ErrMissingSnippet) Error() string {
	return fmt.Sprintf("only metadata in %s", err.path.String())
}

func (err *ErrMissingName) Error() string {
	return "node is missing name"
}

func (err *ErrMalformed) Error() string {
	return fmt.Sprintf("node in %s was malformed", err.path.String())
}
