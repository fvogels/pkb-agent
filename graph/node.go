package graph

type Node interface {
	GetName() string
	GetLinks() []string
}
