package listview

type MsgSetItems[T Item] struct {
	Items List[T]
}

type MsgSelectPrevious struct{}

type MsgSelectNext struct{}

type MsgItemSelected[T Item] struct {
	Index int
	Item  T
}
