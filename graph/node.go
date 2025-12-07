package graph

type Node struct {
	Name      string
	Links     []string
	Backlinks []string
	Extra     any
}
