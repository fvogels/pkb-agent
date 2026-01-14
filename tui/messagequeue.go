package tui

import "github.com/gdamore/tcell/v3"

type MessageQueue interface {
	Enqueue(Message)
}

type messageQueue struct {
	eventQueue chan tcell.Event
}

type EventMessage struct {
	tcell.EventTime
	Message Message
}

func NewMessageQueue(queue chan tcell.Event) MessageQueue {
	return &messageQueue{
		eventQueue: queue,
	}
}

func (queue *messageQueue) Enqueue(message Message) {
	wrapper := EventMessage{
		Message: message,
	}
	wrapper.SetEventNow()

	queue.eventQueue <- &wrapper
}
