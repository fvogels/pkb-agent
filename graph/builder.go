package graph

import (
	"log/slog"
	"pkb-agent/trie"
	"pkb-agent/util"
)

type Builder struct {
	nodes map[string]*Node
}

func NewBuilder() *Builder {
	builder := Builder{
		nodes: make(map[string]*Node),
	}

	return &builder
}

func (builder *Builder) AddNode(node *Node) error {
	if _, alreadyExists := builder.nodes[node.Name]; alreadyExists {
		slog.Debug("Multiple nodes with same name", slog.String("name", node.Name))
		return &ErrNameClash{}
	}

	builder.nodes[node.Name] = node
	return nil
}

func (builder *Builder) Finish() (*Graph, error) {
	if err := builder.ensureLinkedNodeExistence(); err != nil {
		return nil, err
	}

	builder.addBackLinks()

	graph := Graph{
		nodes:    builder.nodes,
		trieRoot: builder.createTrie(),
	}

	return &graph, nil
}

func (builder *Builder) ensureLinkedNodeExistence() error {
	slog.Debug("Checking node links")

	nodes := builder.nodes

	foundUnknownLinks := false
	for _, node := range nodes {
		for _, link := range node.Links {
			if _, found := nodes[link]; !found {
				slog.Debug("Unknown link", slog.String("node", node.Name), slog.String("link target", link))
				foundUnknownLinks = true
			}
		}
	}

	if foundUnknownLinks {
		return &ErrUnknownNodes{}
	}

	return nil
}

func (builder *Builder) addBackLinks() {
	slog.Debug("Adding back links to graph")

	nodes := builder.nodes
	for _, node := range nodes {
		for _, link := range node.Links {
			linkedNode, ok := nodes[link]
			if !ok {
				panic("missing node; should have been noticed earlier")
			}

			linkedNode.Backlinks = append(linkedNode.Backlinks, node.Name)
		}
	}
}

func (builder *Builder) createTrie() *trie.Node[*Node] {
	nodes := builder.nodes
	trieBuilder := trie.NewBuilder[*Node]()

	for name, node := range nodes {
		words := util.LowercaseWords(name)

		for _, word := range words {
			trieBuilder.Add(word, node)
		}
	}

	return trieBuilder.Finish()
}
