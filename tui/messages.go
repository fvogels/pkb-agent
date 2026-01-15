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

type MsgResize struct {
	Size Size
}

type MsgKey struct {
	Key string
}

type MsgUpdateLayout struct{}
