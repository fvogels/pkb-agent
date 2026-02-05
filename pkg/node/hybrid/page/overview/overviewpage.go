package overview

import (
	"pkb-agent/pkg/node"
	"pkb-agent/tui"
)

type Page struct {
}

func New() *Page {
	return &Page{}
}

func (page *Page) CreateViewer(messageQueue tui.MessageQueue) tui.Component {
	return NewPageComponent(messageQueue)
}

func (page *Page) GetCaption() string {
	return "Overview"
}

func (page *Page) GetActions() []node.Action {
	return nil
}
