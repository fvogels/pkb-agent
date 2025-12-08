package listview

type MsgSetItems struct {
	Items List
}

type MsgSelectPrevious struct{}

type MsgSelectNext struct{}

type MsgItemSelected struct {
	Index int
	Item  string
}
