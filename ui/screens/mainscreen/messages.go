package mainscreen

import "pkb-agent/graph"

type MsgGraphLoaded struct {
	graph *graph.Graph
}

type MsgUpdateNodeList struct{}
