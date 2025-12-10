package nodeselectionview

type MsgSetRemainingNodes struct {
	RemainingNodes List
}

type MsgSetSelectedNodes struct {
	SelectedNodes List
}

type MsgSelectNext struct{}

type MsgSelectPrevious struct{}
