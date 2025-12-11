package listview

type MsgSetItems[T any] struct {
	Items List[T]
}

type MsgSelectPrevious struct{}

type MsgSelectNext struct{}

type MsgItemSelected[T any] struct {
	Index int
	Item  T
}

type MsgNoItemSelected struct{}
