package model

import (
	"pkb-agent/persistent/list"
)

type ModelUpdate struct {
	originalModel *Model
	UpdatedModel  Model
}

func (update *ModelUpdate) Apply() {
	original := update.originalModel
	updated := &update.UpdatedModel
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
		update.DetermineIntersectionNodes()
		original.intersectionNodes.Set(updated.intersectionNodes.Get())
	}
}

func (update *ModelUpdate) DetermineIntersectionNodes() {
	original := update.originalModel
	updated := &update.UpdatedModel

	updated.intersectionNodes.Set(determineIntersectionNodes(original.input.Get(), original.Graph, original.selectedNodes.Get(), true, true))
}

func (update *ModelUpdate) SelectHighlightedNode() {
	model := &update.UpdatedModel
	highlightedNode := model.GetHighlightedNode()
	updatedSelectedNodes := list.Append(model.selectedNodes.Get(), highlightedNode)
	model.selectedNodes.Set(updatedSelectedNodes)
}

func (update *ModelUpdate) SetInput(input string) {
	update.UpdatedModel.input.Set(input)
}

func (update *ModelUpdate) Highlight(index int) {
	update.UpdatedModel.highlightedNodeIndex.Set(index)
}

func (update *ModelUpdate) UnselectLastNode() {
	updated := &update.UpdatedModel
	selectedNodes := updated.selectedNodes.Get()

	if selectedNodes.Size() > 0 {
		updated.selectedNodes.Set(list.DropLast(selectedNodes))
	}
}
