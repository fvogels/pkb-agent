package nodeselection

import (
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/component/docknorth"
	"pkb-agent/tui/component/stringlist"
	"pkb-agent/tui/component/stringsview"
	"pkb-agent/tui/data"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	size                 tui.Size
	selectedNodes        data.List[*pkg.Node]
	nodeIntersection     data.List[*pkg.Node]
	selectedIndex        data.Value[int]
	selectedNodesView    *stringsview.Component
	nodeIntersectionView *stringlist.Component
	root                 *docknorth.Component
}

func New(selectedNodes data.List[*pkg.Node], nodeIntersection data.List[*pkg.Node], selectedIndex data.Value[int]) *Component {
	style := tcell.StyleDefault.Background(color.Green)
	selectedNodesNames := data.MapList(selectedNodes, func(node *pkg.Node) stringsview.Item {
		name := node.GetName()
		item := stringsview.Item{
			Runes: []rune(name),
			Style: &style,
		}
		return item
	})

	selectedNodesView := stringsview.New(selectedNodesNames)

	nodeIntersectionItems := data.MapList(nodeIntersection, func(node *pkg.Node) string {
		return node.GetName()
	})
	nodeIntersectionView := stringlist.New(nodeIntersectionItems, selectedIndex)

	root := docknorth.New(
		"nodeselection[selected|intersection]",
		selectedNodesView,
		nodeIntersectionView,
		0,
	)

	component := Component{
		selectedNodes:        selectedNodes,
		nodeIntersection:     nodeIntersection,
		selectedIndex:        selectedIndex,
		selectedNodesView:    selectedNodesView,
		nodeIntersectionView: nodeIntersectionView,
		root:                 root,
	}

	component.updateLayout()
	selectedNodes.Observe(func() { component.updateLayout() })

	return &component
}

func (component *Component) SetOnSelectionChanged(callback func(int)) {
	component.nodeIntersectionView.SetOnSelectionChanged(callback)
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	default:
		component.root.Handle(message)
	}
}

func (component *Component) Render() tui.Grid {
	return component.root.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size
	component.root.Handle(message)
}

func (component *Component) updateLayout() {
	selectedNodeCount := component.selectedNodes.Size()
	component.root.SetDockerChildHeight(selectedNodeCount)
}
