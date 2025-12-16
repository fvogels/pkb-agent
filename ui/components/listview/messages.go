package listview

type MsgSetItems[T any] struct {
	Items          List[T]
	SelectionIndex int
}

type MsgSelectItem struct {
	Index int
}

type MsgItemSelected[T any] struct {
	Index int
	Item  T
}

type MsgNoItemSelected struct{}

type MsgSetItemRenderer[T any] struct {
	Renderer func(item T) string
}
