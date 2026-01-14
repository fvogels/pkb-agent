package tagger

import "pkb-agent/pkg/node"

type NodeWrapper struct {
	node.RawNode
	extraLinks []string
}

func (wrapper *NodeWrapper) GetLinks() []string {
	return append(node.RawNode.GetLinks(wrapper.RawNode), wrapper.extraLinks...)
}
