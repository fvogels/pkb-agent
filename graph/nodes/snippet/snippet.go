package snippet

import pathlib "pkb-agent/util/path"

type Node struct {
	Name       string
	Identifier string
	Links      []string
	Path       pathlib.Path
}

func (node *Node) GetName() string {
	return node.Name
}

func (node *Node) GetLinks() []string {
	return node.Links
}
