package graph

type Node struct {
	Index     int
	Name      string
	Type      string
	Links     []string
	Backlinks []string
	Info      any
}
