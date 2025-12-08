package trie

type Node[T any] struct {
	Children     []*Node[T]
	Terminals    []T
	NextTerminal *Node[T]
}

func (node *Node[T]) walk(callback func(*Node[T])) {
	callback(node)

	if node.Children != nil {
		for _, child := range node.Children {
			if child != nil {
				child.walk(callback)
			}
		}
	}
}
