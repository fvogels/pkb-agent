package model

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
)

type Model struct {
	Graph                *pkg.Graph
	SelectedNodes        list.List[*pkg.Node]
	IntersectionNodes    list.List[*pkg.Node]
	HighlightedNodeIndex int
	Input                string
	ShowNodeLinks        bool // Whether to show node contents (false) or node links/backlinks (true)
}

func New(graph *pkg.Graph) *Model {
	selectedNodes := list.New[*pkg.Node]()
	intersectionNodes := list.FromSlice(graph.ListNodes())
	highlightedNodeIndex := 0

	result := Model{
		Graph:                graph,
		SelectedNodes:        selectedNodes,
		IntersectionNodes:    intersectionNodes,
		HighlightedNodeIndex: highlightedNodeIndex,
		ShowNodeLinks:        false,
	}

	return &result
}

func (model *Model) GetHighlightedNode() *pkg.Node {
	return model.IntersectionNodes.At(model.HighlightedNodeIndex)
}

func (model *Model) DetermineIntersectionNodes() {
	model.IntersectionNodes = determineIntersectionNodes(model.Input, model.Graph, model.SelectedNodes, true, true)
}

func (model *Model) SelectHighlightedNode() {
	highlightedNode := model.GetHighlightedNode()
	updatedSelectedNodes := list.Append(model.SelectedNodes, highlightedNode)
	model.SelectedNodes = updatedSelectedNodes
}

func (model *Model) UnselectLastNode() {
	selectedNodes := model.SelectedNodes

	if selectedNodes.Size() > 0 {
		model.SelectedNodes = list.DropLast(selectedNodes)
	}
}
