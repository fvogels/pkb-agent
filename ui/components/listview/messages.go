package listview

type MsgSetItems[T any] struct {
	Items          List[T]
	SelectionIndex int
}

type MsgSelectPrevious struct{}

type MsgSelectNext struct{}

type MsgItemSelected[T any] struct {
	Index int
	Item  T
}

type MsgNoItemSelected struct{}

type MsgSetItemRenderer[T any] struct {
	Renderer func(item T) string
}
