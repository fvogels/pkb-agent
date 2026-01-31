package nodeselection

import (
	"fmt"
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/component/border"
	"pkb-agent/tui/component/docknorth"
	"pkb-agent/tui/component/stringlist"
	"pkb-agent/tui/component/stringsview"
	"pkb-agent/tui/data"
	"pkb-agent/tui/grid"
	"pkb-agent/tui/position"
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
	root                 tui.Component
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

	borderStyle := tcell.StyleDefault.Foreground(color.Reset).Background(color.Reset)
	root := border.New(
		messageQueue,
		docknorth.New(
			messageQueue,
			"nodeselection[selected|intersection]",
			selectedNodesView,
			nodeIntersectionView,
			data.MapValue(selectedNodes, func(nodes list.List[*pkg.Node]) int { return nodes.Size() }),
		),
		&borderStyle,
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

func (component *Component) Render() grid.FiniteGrid {
	result := grid.Materialize(component.root.Render())

	component.renderNodeCount(result)

	return result
}

func (component *Component) renderNodeCount(g *grid.MaterializedGrid) {
	str := fmt.Sprintf(" %d/%d ", component.selectedIndex.Get()+1, component.nodeIntersection.Get().Size())

	// Only add node count if the grid is wide enough
	if len(str)+4 < g.Size().Width {
		style := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
		x := g.Size().Width - 2
		y := g.Size().Height - 1

		for index, char := range str {
			pos := position.Position{
				X: x - len(str) + index,
				Y: y,
			}

			cell := grid.Cell{
				Contents: char,
				Style:    &style,
			}

			g.Set(pos, cell)
		}
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.root.Handle(message)
}

func (component *Component) updateLayout() {
	component.root.Handle(tui.MsgResize{
		Size: component.Size,
	})
}

func (component *Component) onStateUpdated() {
	component.updateLayout()
}
