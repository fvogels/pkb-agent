package tui

import "fmt"

type Message interface {
	String() string
}

const (
	Everyone = -1
)

type MsgActivate struct {
	Recipient int
}

func (message MsgActivate) String() string {
	return fmt.Sprintf("MsgActivate[Recipient=%d]", message.Recipient)
}

func (message MsgActivate) ShouldRespond(identifier int) bool {
	return message.Recipient == Everyone || message.Recipient == identifier
}

func (message MsgActivate) Respond(receiver int, onActivate func(), children ...Component) {
	if message.ShouldRespond(receiver) {
		onActivate()

		for _, child := range children {
			child.Handle(MsgActivate{Recipient: Everyone})
		}
	} else {
		for _, child := range children {
			child.Handle(message)
		}
	}
}

type MsgResize struct {
	Size Size
}

func (message MsgResize) String() string {
	return fmt.Sprintf("MsgResize[Size=%s]", message.Size.String())
}

type MsgKey struct {
	Key string
}

func (message MsgKey) String() string {
	return fmt.Sprintf("MsgKeyKey=%s]", message.Key)
}

type MsgUpdateLayout struct{}

func (message MsgUpdateLayout) String() string {
	return fmt.Sprintf("MsgCommand[...]")
}

type MsgCommand struct {
	Command func()
}

func (message MsgCommand) String() string {
	return "MsgCommand[...]"
}
