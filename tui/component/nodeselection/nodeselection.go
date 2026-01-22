package nodeselection

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/component/docknorth"
	"pkb-agent/tui/component/stringlist"
	"pkb-agent/tui/component/stringsview"
	"pkb-agent/tui/data"
	"pkb-agent/tui/grid"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	lockCount            data.Value[int]
	selectedNodes        data.Value[list.List[*pkg.Node]]
	nodeIntersection     data.Value[list.List[*pkg.Node]]
	selectedIndex        data.Value[int]
	selectedNodesView    *stringsview.Component
	nodeIntersectionView *stringlist.Component
	root                 *docknorth.Component
}

func New(messageQueue tui.MessageQueue, selectedNodes data.Value[list.List[*pkg.Node]], nodeIntersection data.Value[list.List[*pkg.Node]], selectedIndex data.Value[int], lockCount data.Value[int]) *Component {
	lockedStyle := tcell.StyleDefault.Background(color.Red)
	unlockedStyle := tcell.StyleDefault.Background(color.Green)

	selectedNodesNames := data.Cache(
		data.MapValue2(selectedNodes, lockCount, func(selectedNodes list.List[*pkg.Node], lockCount int) list.List[stringsview.Item] {
			return list.MapWithIndex(selectedNodes, func(index int, node *pkg.Node) stringsview.Item {
				name := node.GetName()

				var style *tui.Style
				if index < lockCount {
					style = &lockedStyle
				} else {
					style = &unlockedStyle
				}

				item := stringsview.Item{
					Runes: []rune(name),
					Style: style,
				}
				return item
			})
		}),
	)

	selectedNodesView := stringsview.New(messageQueue, selectedNodesNames)

	nodeIntersectionItems := data.Cache(
		data.MapValue(nodeIntersection, func(lst list.List[*pkg.Node]) list.List[string] {
			return list.MapList(lst, func(node *pkg.Node) string {
				return node.GetName()
			})
		}),
	)
	nodeIntersectionView := stringlist.New(messageQueue, nodeIntersectionItems, selectedIndex)

	root := docknorth.New(
		messageQueue,
		"nodeselection[selected|intersection]",
		selectedNodesView,
		nodeIntersectionView,
		0,
	)

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed node selection view",
			MessageQueue: messageQueue,
		},
		selectedNodes:        selectedNodes,
		nodeIntersection:     nodeIntersection,
		selectedIndex:        selectedIndex,
		selectedNodesView:    selectedNodesView,
		nodeIntersectionView: nodeIntersectionView,
		root:                 root,
	}

	component.updateLayout()

	return &component
}

func (component *Component) SetOnSelectionChanged(callback func(int)) {
	component.nodeIntersectionView.SetOnSelectionChanged(callback)
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	case tui.MsgStateUpdated:
		component.onStateUpdated()

	default:
		component.root.Handle(message)
	}
}

func (component *Component) Render() grid.Grid {
	return component.root.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.root.Handle(message)
}

func (component *Component) updateLayout() {
	selectedNodeCount := component.selectedNodes.Get().Size()
	component.root.SetDockerChildHeight(selectedNodeCount)
}

func (component *Component) onStateUpdated() {
	component.updateLayout()
}
