package model

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui/data"
)

type Model struct {
	Graph                *pkg.Graph
	selectedNodes        data.Variable[list.List[*pkg.Node]]
	intersectionNodes    data.Variable[list.List[*pkg.Node]]
	highlightedNodeIndex data.Variable[int]
	input                data.Variable[string]
}

func New(graph *pkg.Graph) Model {
	selectedNodes := list.New[*pkg.Node]()
	intersectionNodes := list.FromSlice(graph.ListNodes())
	highlightedNodeIndex := 0

	result := Model{
		Graph:                graph,
		selectedNodes:        data.NewVariable(selectedNodes),
		intersectionNodes:    data.NewVariable(intersectionNodes),
		highlightedNodeIndex: data.NewVariable(highlightedNodeIndex),
	}

	return result
}

func (model *Model) SelectedNodes() data.Value[list.List[*pkg.Node]] {
	return &model.selectedNodes
}

func (model *Model) IntersectionNodes() data.Value[list.List[*pkg.Node]] {
	return &model.intersectionNodes
}

func (model *Model) Input() data.Value[string] {
	return &model.input
}

func (model *Model) HighlightedNodeIndex() data.Value[int] {
	return &model.highlightedNodeIndex
}

func (model *Model) Update() *ModelUpdate {
	return &ModelUpdate{
		originalModel: model,
		UpdatedModel:  *model,
	}
}

func (model *Model) GetHighlightedNode() *pkg.Node {
	return model.intersectionNodes.Get().At(model.highlightedNodeIndex.Get())
}
