package application

import (
	"fmt"
	"log/slog"
	"os"
	"pkb-agent/graph/loaders/sequence"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/component/stringsview"
	"pkb-agent/tui/data"
	"pkb-agent/util/pathlib"
	"slices"
	"strings"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

const (
	logFilename = "ui.log"
)

type Application struct {
	verbose    bool
	running    bool
	logFile    *os.File
	screen     tcell.Screen
	size       tui.Size
	graph      *pkg.Graph
	model      Model
	viewMode   *viewMode
	inputMode  *inputMode
	activeMode mode
}

type Model struct {
	input                 *data.Variable[string]
	selectedNodes         *data.SliceList[*pkg.Node]
	intersectionNodes     *data.SliceList[*pkg.Node]
	highlightedNodeIndex  *data.Variable[int]
	intersectionNodeNames data.List[string]
}

func NewApplication(verbose bool) *Application {
	application := Application{
		verbose: verbose,
		running: true,
	}

	return &application
}

func (application *Application) Start() error {
	err := application.initializeLogging()
	if err != nil {
		return err
	}

	if err := application.initializeScreen(); err != nil {
		return err
	}

	if err := application.loadGraph(); err != nil {
		return err
	}

	application.createModel()
	application.viewMode = newViewMode(application)
	application.inputMode = newInputMode(application)
	application.activeMode = application.viewMode

	application.eventLoop()

	return nil
}

func (application *Application) Close() {
	maybePanic := recover()
	application.screen.Fini()
	if maybePanic != nil {
		panic(maybePanic)
	}

	if application.logFile != nil {
		application.logFile.Close()
	}
}

func (application *Application) initializeScreen() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	if err := screen.Init(); err != nil {
		return err
	}

	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	screen.SetStyle(defStyle)
	screen.EnableMouse()
	// s.EnablePaste()
	screen.Clear()

	application.screen = screen

	return nil
}

func (application *Application) eventLoop() {
	// style := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	screen := application.screen

	for application.running {
		activeMode := application.activeMode

		application.Render()

		// Poll event (this can be in a select statement as well)
		ev := <-screen.EventQ()

		// Process event
		switch event := ev.(type) {
		case *tcell.EventResize:
			width, height := event.Size()

			application.size = tui.Size{
				Width:  width,
				Height: height,
			}

			activeMode.Handle(tui.MsgResize{
				Size: application.size,
			})
			screen.Sync()

		case *tcell.EventKey:
			translation := translateKey(event)
			slog.Debug("Key pressed", slog.String("key", translation))

			message := tui.MsgKey{
				Key: translateKey(event),
			}
			activeMode.Handle(message)

			// case *tcell.EventMouse:
			// 	x, y := event.Position()
			// 	position := tui.Position{X: x, Y: y}
			// 	clickHandler := grid.Get(position).OnClick

			// 	if clickHandler != nil && event.Buttons() == tcell.Button1 {
			// 		clickHandler()
			// 	}
		}
	}
}

func (application *Application) Render() {
	screen := application.screen
	activeMode := application.activeMode

	screen.Clear()

	grid := activeMode.Render()
	gridSize := grid.GetSize()
	runes := make([]rune, 1)

	timeBeforeUpdate := time.Now()

	for y := range gridSize.Height {
		for x := range gridSize.Width {
			position := tui.Position{X: x, Y: y}
			cell := grid.Get(position)
			runes[0] = cell.Contents
			screen.Put(x, y, string(runes), *cell.Style)
		}
	}

	screen.Show()
	slog.Debug("Screen updated", slog.String("duration", time.Since(timeBeforeUpdate).String()))
}

func (application *Application) createModel() {
	graph := application.graph

	// Create data sources
	input := data.NewVariable("")
	highlightedNodeIndex := data.NewVariable(0)
	selectedNodes := data.NewSliceList[*pkg.Node](nil)
	intersectionNodes := data.NewSliceList[*pkg.Node](nil)
	intersectionNodeNames := data.MapList(intersectionNodes, func(node *pkg.Node) string { return node.GetName() })

	// Cause intersection node list to be updated whenever the input or the selected nodes change
	updateIntersectionNodes := func() {
		nodes := determineIntersectionNodes(input.Get(), graph, data.CopyListToSlice(selectedNodes), true, true)
		intersectionNodes.SetSlice(nodes)
	}
	updateIntersectionNodes()
	data.DefineReaction(updateIntersectionNodes, input, selectedNodes)

	input.Observe(func() {
		if len(input.Get()) > 0 {
			application.updateIntersectionNodeSelection(input.Get())
		}
	})

	model := Model{
		input:                 input,
		selectedNodes:         selectedNodes,
		intersectionNodes:     intersectionNodes,
		highlightedNodeIndex:  highlightedNodeIndex,
		intersectionNodeNames: intersectionNodeNames,
	}

	application.model = model
}

func (application *Application) initializeLogging() error {
	if application.verbose {
		logFile, err := os.Create(logFilename)
		if err != nil {
			fmt.Println("Failed to create log")
			return err
		}
		application.logFile = logFile

		logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)

		return nil
	}

	return nil
}

func translateKey(event *tcell.EventKey) string {
	if event.Key() == tcell.KeyRune {
		return event.Str()
	} else {
		return tcell.KeyNames[event.Key()]
	}
}

type ItemList struct {
	items [][]rune
	style *tui.Style
}

func (list ItemList) Size() int {
	return len(list.items)
}

func (list ItemList) At(index int) stringsview.Item {
	return stringsview.Item{
		Runes: list.items[index],
		Style: list.style,
	}
}

func (application *Application) loadGraph() error {
	before := time.Now()
	loader := sequence.New()
	path := pathlib.New(`F:\repos\pkb\pkb-data\root.yaml`)

	graph, err := pkg.LoadGraph(path, loader)
	if err != nil {
		return err
	}
	slog.Debug("Graph loaded", slog.String("loadTime", time.Since(before).String()))

	application.graph = graph

	return nil
}

func (application *Application) updateIntersectionNodeSelection(target string) {
	nodes := data.CopyListToSlice(application.model.intersectionNodes)

	bestMatchIndex, found := slices.BinarySearchFunc(
		nodes,
		target,
		func(node *pkg.Node, target string) int {
			nodeName := strings.ToLower(node.GetName())
			if strings.HasPrefix(nodeName, target) {
				return 0
			}
			if nodeName < target {
				return -1
			}
			return 1
		},
	)

	if !found {
		bestMatchIndex = 0
	}

	application.model.highlightedNodeIndex.Set(bestMatchIndex)
}

func (application *Application) selectHighlightedNode() {
	if application.model.intersectionNodes.Size() > 0 {
		highlightedNodeIndex := application.model.highlightedNodeIndex.Get()
		highlightedNode := application.model.intersectionNodes.At(highlightedNodeIndex)
		application.model.selectedNodes.Update(func(ns []*pkg.Node) []*pkg.Node {
			return append(ns, highlightedNode)
		})
	}
}

func (application *Application) unselectLastNode() {
	if application.model.selectedNodes.Size() > 0 {
		application.model.selectedNodes.Update(func(ns []*pkg.Node) []*pkg.Node {
			return ns[:len(ns)-1]
		})
	}
}

func (application *Application) switchMode(mode mode) {
	application.activeMode = mode
	application.activeMode.Handle(tui.MsgResize{Size: application.size})
}
