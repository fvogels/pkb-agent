package trie

type Node[T any] struct {
	Depth             int
	Children          []*Node[T]
	Terminals         []T
	NextTerminal      *Node[T]
	NextTerminalDepth int
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

func (node *Node[T]) walkBackwards(callback func(*Node[T])) {
	if node.Children != nil {
		for i := len(node.Children) - 1; i >= 0; i-- {
			child := node.Children[i]
			if child != nil {
				child.walkBackwards(callback)
			}
		}
	}

	callback(node)
}

// Descend finds the descendant one reaches by following each character in the given string.
// If no such descendant exists, nil is returned.
func (node *Node[T]) Descend(str string) *Node[T] {
	current := node

	for _, char := range str {
		if current.Children == nil {
			return nil
		}

		if current.Children[char] == nil {
			return nil
		}

		current = current.Children[char]
	}

	return current
}
