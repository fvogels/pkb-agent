package pkg

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"pkb-agent/graph/node"
	"pkb-agent/util/trie"
	"slices"
	"strings"
)

type Builder struct {
	rawNodes map[string]node.RawNode
}

func NewBuilder() *Builder {
	builder := Builder{
		rawNodes: make(map[string]node.RawNode),
	}

	return &builder
}

func (builder *Builder) AddNode(node node.RawNode) error {
	if _, alreadyExists := builder.rawNodes[node.GetName()]; alreadyExists {
		return fmt.Errorf("%w: %s", ErrNameClash, node.GetName())
	}

	builder.rawNodes[node.GetName()] = node
	return nil
}

func (builder *Builder) Finalize() (*Graph, error) {
	slog.Debug("Finalizing graph")

	slog.Debug("Wrapping nodes")
	nodes := builder.wrapNodes()

	slog.Debug("Building nodesByName table")
	nodesByName, err := builder.buildNameToNodeTable(nodes)
	if err != nil {
		return nil, err
	}

	slog.Debug("Sorting nodes by name")
	builder.sortByName(nodes)

	slog.Debug("Indexing nodes")
	builder.addIndices(nodes)

	slog.Debug("Linking nodes")
	if err := builder.linkNodes(nodesByName); err != nil {
		return nil, err
	}

	slog.Debug("Backlinking nodes")
	builder.addBackLinks(nodes)

	slog.Debug("Creating trie")
	trie := builder.createTrie(nodes)

	graph := Graph{
		nodesByIndex: nodes,
		nodesByName:  nodesByName,
		trieRoot:     trie,
	}

	slog.Debug("Finished finalizing graph")

	return &graph, nil
}

func (builder *Builder) wrapNodes() []*Node {
	result := []*Node{}

	maps.Values(builder.rawNodes)(func(rawNode node.RawNode) bool {
		wrapper := Node{
			rawNode: rawNode,
		}

		result = append(result, &wrapper)

		return true
	})

	return result
}

func (builder *Builder) sortByName(wrappers []*Node) {
	slices.SortFunc(wrappers, func(first *Node, second *Node) int {
		firstName := strings.ToLower(first.rawNode.GetName())
		secondName := strings.ToLower(second.rawNode.GetName())

		return strings.Compare(firstName, secondName)
	})
}

func (builder *Builder) addIndices(wrappers []*Node) {
	for index, wrapper := range wrappers {
		wrapper.id = index
	}
}

func (builder *Builder) buildNameToNodeTable(wrappers []*Node) (map[string]*Node, error) {
	result := make(map[string]*Node)
	errs := []error{}

	for _, node := range wrappers {
		nodeName := node.rawNode.GetName()

		if _, found := result[nodeName]; found {
			errs = append(errs, fmt.Errorf("%w: multiple nodes with name \"%s\"", ErrNameClash, nodeName))
		} else {
			result[nodeName] = node
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return result, nil
}

func (builder *Builder) linkNodes(table map[string]*Node) error {
	errs := []error{}

	maps.Values(table)(func(node *Node) bool {
		for _, linkedName := range node.rawNode.GetLinks() {
			linkedNode, found := table[linkedName]

			if !found {
				err := fmt.Errorf("%w: \"%s\" links to unknown node \"%s\"", ErrUnknownNode, node.rawNode.GetName(), linkedName)
				errs = append(errs, err)
			}

			node.links = append(node.links, linkedNode)
		}

		return true
	})

	return errors.Join(errs...)
}

func (builder *Builder) addBackLinks(nodes []*Node) {
	for _, node := range nodes {
		for _, linkedNode := range node.links {
			linkedNode.backlinks = append(linkedNode.backlinks, node)
		}
	}
}

func (builder *Builder) createTrie(nodes []*Node) *trie.Node[*Node] {
	trieBuilder := trie.NewBuilder[*Node]()

	for _, node := range nodes {
		for _, keyword := range node.rawNode.GetSearchStrings() {
			trieBuilder.Add(keyword, node)
		}
	}

	return trieBuilder.Finish()
}
