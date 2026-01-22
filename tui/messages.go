package tui

import (
	"fmt"
	"pkb-agent/tui/size"
)

type Message interface {
	String() string
}

const (
	Everyone = -1
)

type MsgResize struct {
	Size size.Size
}

func (message MsgResize) String() string {
	return fmt.Sprintf("MsgResize[Size=%s]", message.Size.String())
}

type MsgKey struct {
	Key string
}

func (message MsgKey) String() string {
	return fmt.Sprintf("MsgKey[Key=%s]", message.Key)
}

type MsgUpdateLayout struct{}

func (message MsgUpdateLayout) String() string {
	return "MsgCommand"
}

type MsgCommand struct {
	Command func()
}

func (message MsgCommand) String() string {
	return "MsgCommand[...]"
}

type MsgStateUpdated struct{}

func (message MsgStateUpdated) String() string {
	return "MsgStateUpdated"
}
