package model

import (
	"pkb-agent/persistent/list"
)

type ModelUpdate struct {
	originalModel *Model
	updatedModel  Model
}

func (update *ModelUpdate) Apply() {
	original := update.originalModel
	updated := &update.updatedModel
	updateIntersection := false

	if original.selectedNodes.Get() != updated.selectedNodes.Get() {
		original.selectedNodes.Set(updated.selectedNodes.Get())
		updateIntersection = true
	}

	if original.input.Get() != updated.input.Get() {
		original.input.Set(updated.input.Get())
		updateIntersection = true
	}

	if original.highlightedNodeIndex.Get() != updated.highlightedNodeIndex.Get() {
		original.highlightedNodeIndex.Set(updated.highlightedNodeIndex.Get())
	}

	if updateIntersection {
		updatedIntersectionNodes := determineIntersectionNodes(original.input.Get(), original.Graph, original.selectedNodes.Get(), true, true)
		original.intersectionNodes.Set(updatedIntersectionNodes)
	}
}

func (update *ModelUpdate) SelectHighlightedNode() {
	model := &update.updatedModel
	highlightedNode := model.GetHighlightedNode()
	updatedSelectedNodes := list.Append(model.selectedNodes.Get(), highlightedNode)
	model.selectedNodes.Set(updatedSelectedNodes)
}

func (update *ModelUpdate) SetInput(input string) {
	update.updatedModel.input.Set(input)
}

func (update *ModelUpdate) Highlight(index int) {
	update.updatedModel.highlightedNodeIndex.Set(index)
}

func (update *ModelUpdate) UnselectLastNode() {
	updated := &update.updatedModel
	selectedNodes := updated.selectedNodes.Get()

	if selectedNodes.Size() > 0 {
		updated.selectedNodes.Set(list.DropLast(selectedNodes))
	}
}
