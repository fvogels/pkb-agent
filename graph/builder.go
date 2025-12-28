package graph

import (
	"log/slog"
	"maps"
	"pkb-agent/util"
	"pkb-agent/util/trie"
	"slices"
	"strings"
	"unicode"
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
		slog.Error("Multiple nodes with same name", slog.String("name", node.Name))
		return &ErrNameClash{
			name: node.Name,
		}
	}

	builder.nodes[node.Name] = node
	return nil
}

func (builder *Builder) Finish() (*Graph, error) {
	if err := builder.ensureLinkedNodeExistence(); err != nil {
		return nil, err
	}

	nodesByIndex := builder.addIndices()
	builder.addBackLinks()

	graph := Graph{
		nodesByIndex: nodesByIndex,
		nodesByName:  builder.nodes,
		trieRoot:     builder.createTrie(),
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
				slog.Error("Unknown link", slog.String("node", node.Name), slog.String("link target", link))
				foundUnknownLinks = true
			}
		}
	}

	if foundUnknownLinks {
		return &ErrUnknownNodes{}
	}

	return nil
}

func (builder *Builder) addIndices() []*Node {
	nodes := []*Node{}
	maps.Values(builder.nodes)(func(node *Node) bool {
		nodes = append(nodes, node)
		return true
	})

	slices.SortFunc(nodes, func(a, b *Node) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	for index, node := range nodes {
		node.Index = index
	}

	return nodes
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
		searchPrefixes := builder.deriveSearchPrefixes(name)

		for _, searchPrefix := range searchPrefixes {
			trieBuilder.Add(searchPrefix, node)
		}
	}

	return trieBuilder.Finish()
}

func (builder *Builder) deriveSearchPrefixes(nodeName string) []string {
	s := nodeName
	s = strings.TrimSpace(s)
	s = util.RemoveAccents(s)
	s = strings.ToLower(s)
	s = util.KeepOnlyLettersAndSpaces(s)

	prefixes := []string{s}

	for index, rune := range s {
		if unicode.IsSpace(rune) {
			prefixes = append(prefixes, s[index+1:])
		}
	}

	return prefixes
}
