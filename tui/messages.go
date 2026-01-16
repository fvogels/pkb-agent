package tui

type Message any

const (
	Everyone = -1
)

type MsgActivate struct {
	Recipient int
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

type MsgKey struct {
	Key string
}

type MsgUpdateLayout struct{}
