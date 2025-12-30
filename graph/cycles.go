package graph

import "log/slog"

func ContainsCycles(graph *Graph) bool {
	nodeCount := graph.GetNodeCount()

	detector := cycleDetector{
		graph:      graph,
		visited:    make([]bool, nodeCount),
		deemedSafe: make([]bool, nodeCount),
	}

	for i := range nodeCount {
		slog.Debug("Looking for cycle", slog.String("startNode", graph.FindNodeByIndex(i).rawNode.GetName()))

		if detector.detectCycles(i) {
			node := graph.FindNodeByIndex(i)

			slog.Error("Cycle detected", slog.String("nodeName", node.rawNode.GetName()))
			return true
		}
	}

	return false
}

type cycleDetector struct {
	graph      *Graph
	visited    []bool
	deemedSafe []bool
}

func (detector *cycleDetector) detectCycles(nodeIndex int) bool {
	if detector.deemedSafe[nodeIndex] {
		// Used cached result
		return false
	}

	if detector.visited[nodeIndex] {
		// We encountered this node earlier, meaning there's a cycle
		return true
	}

	graph := detector.graph
	detector.visited[nodeIndex] = true
	node := graph.FindNodeByIndex(nodeIndex)

	for _, linked := range node.links {
		if detector.detectCycles(linked.id) {
			return true
		}
	}

	detector.deemedSafe[nodeIndex] = true
	return false
}
