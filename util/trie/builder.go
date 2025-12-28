package trie

type Builder[T any] struct {
	root *Node[T]
}

func NewBuilder[T any]() *Builder[T] {
	root := Node[T]{
		Depth:             0,
		Children:          nil,
		Terminals:         nil,
		NextTerminal:      nil,
		NextTerminalDepth: 0,
	}

	builder := Builder[T]{
		root: &root,
	}

	return &builder
}

func (builder *Builder[T]) Add(string string, terminal T) {
	current := builder.root

	for _, char := range string {
		if current.Children == nil {
			current.Children = make([]*Node[T], 128)
		}
		if current.Children[char] == nil {
			current.Children[char] = &Node[T]{
				Depth: current.Depth + 1,
			}
		}
		current = current.Children[char]
	}

	if current.Terminals == nil {
		current.Terminals = make([]T, 0, 1)
	}

	current.Terminals = append(current.Terminals, terminal)
}

func (builder *Builder[T]) addLinks() {
	var lastTerminal *Node[T] = nil
	minDepth := 0

	walker := func(node *Node[T]) {
		node.NextTerminal = lastTerminal
		node.NextTerminalDepth = minDepth - 1

		if node.Depth < minDepth {
			minDepth = node.Depth
		}

		if len(node.Terminals) > 0 {
			lastTerminal = node
			minDepth = node.Depth
		}
	}

	builder.root.walkBackwards(walker)
}

func (builder *Builder[T]) Finish() *Node[T] {
	builder.addLinks()
	return builder.root
}
