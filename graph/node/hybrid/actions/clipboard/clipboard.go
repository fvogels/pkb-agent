package clipboard

import (
	"golang.design/x/clipboard"
)

type Action struct {
	contents string
}

func New(contents string) *Action {
	return &Action{
		contents: contents,
	}
}

func (action Action) GetDescription() string {
	return "copy"
}

func (action Action) Perform() {
	clipboard.Write(clipboard.FmtText, ([]byte)(action.contents))
}
