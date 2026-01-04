package hybrid

import "errors"

var (
	ErrMissingMetadata   = errors.New("missing metadata section")
	ErrMissingName       = errors.New("missing name")
	ErrMissingActionType = errors.New("missing type in action definition")
	ErrUnknownActionType = errors.New("unknown action type")
)
