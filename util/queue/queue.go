package queue

type Queue[T any] struct {
	items []T
}

func New[T any]() *Queue[T] {
	return &Queue[T]{items: nil}
}

func (queue *Queue[T]) Enqueue(item T) {
	queue.items = append(queue.items, item)
}

func (queue *Queue[T]) IsEmpty() bool {
	return len(queue.items) == 0
}

func (queue *Queue[T]) Dequeue() T {
	result := queue.items[len(queue.items)-1]
	queue.items = queue.items[:len(queue.items)-1]
	return result
}
