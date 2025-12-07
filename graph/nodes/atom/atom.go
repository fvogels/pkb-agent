package atom

type Node struct {
	Name       string
	Identifier string
	Links      []string
}

func (node *Node) GetName() string {
	return node.Name
}

func (node *Node) GetLinks() []string {
	return node.Links
}
