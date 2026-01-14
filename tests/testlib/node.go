package testlib

import (
	"io"
	"pkb-agent/tui"
)

type TestNode struct {
	Name     string
	Keywords []string
	Links    []string
}

func (node *TestNode) GetName() string {
	return node.Name
}

func (node *TestNode) GetSearchStrings() []string {
	return node.Keywords
}

func (node *TestNode) GetLinks() []string {
	return node.Links
}

func (node *TestNode) GetViewer(tui.MessageQueue) tui.Component {
	return nil
}

func (node *TestNode) Serialize(io.Writer) error {
	return nil
}
