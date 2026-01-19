package tui

import (
	"log/slog"
	"runtime"

	"github.com/gdamore/tcell/v3"
)

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
	{
		_, file, line, _ := runtime.Caller(1)
		slog.Debug(
			"Message enqueued",
			slog.String("file", file),
			slog.Int("line", line),
			slog.String("message", message.String()),
		)
	}

	wrapper := EventMessage{
		Message: message,
	}
	wrapper.SetEventNow()

	queue.eventQueue <- &wrapper
}
