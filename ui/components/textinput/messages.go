package textinput

type MsgInputUpdated struct {
	Input string
}

type MsgClear struct{}

type MsgSetInput struct {
	Input string
}
