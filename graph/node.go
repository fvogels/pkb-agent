package graph

type Node struct {
	Name      string
	Type      string
	Links     []string
	Backlinks []string
	Info      any
}
