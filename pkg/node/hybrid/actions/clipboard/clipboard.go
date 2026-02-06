package clipboard

import (
	"golang.design/x/clipboard"
)

type Action struct {
	contents string
	key      string
}

func New(contents string, key string) *Action {
	return &Action{
		contents: contents,
		key:      key,
	}
}

func (action Action) GetDescription() string {
	return "copy"
}

func (action Action) GetKey() string {
	return action.key
}

func (action Action) Perform() {
	clipboard.Write(clipboard.FmtText, ([]byte)(action.contents))
}
