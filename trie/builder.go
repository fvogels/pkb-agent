package trie

type Builder[T any] struct {
	root *Node[T]
}

func (builder *Builder[T]) Add(string string, terminal T) {
	current := builder.root

	for _, char := range string {
		if current.Children == nil {
			current.Children = make([]*Node[T], 128)
		}
		if current.Children[char] == nil {
			current.Children[char] = &Node[T]{}
		}
		current = current.Children[char]
	}

	if current.Terminals == nil {
		current.Terminals = make([]T, 0, 1)
	}

	current.Terminals = append(current.Terminals, terminal)
}

func (builder *Builder[T]) AddLinks() {
	var lastTerminal *Node[T] = nil

	walker := func(node *Node[T]) {
		if node.Terminals != nil {
			if lastTerminal != nil {
				lastTerminal.NextTerminal = node
			}
			lastTerminal = node
		}
	}

	builder.root.walk(walker)
}
