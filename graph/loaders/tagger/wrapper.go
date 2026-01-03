package tagger

import "pkb-agent/graph/node"

type NodeWrapper struct {
	node.RawNode
	extraLinks []string
}

func (wrapper *NodeWrapper) GetLinks() []string {
	return append(node.RawNode.GetLinks(wrapper.RawNode), wrapper.extraLinks...)
}
