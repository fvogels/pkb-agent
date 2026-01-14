package application

import (
	"fmt"
	"log/slog"
	"os"
	"pkb-agent/graph/loaders/sequence"
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/model"
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
	model      model.Model
	viewMode   *viewMode
	inputMode  *inputMode
	activeMode mode
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

	application.model = model.New(application.graph)
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
	for application.running {
		application.Render()

		event := application.WaitForEvent()
		for event != nil && application.running {
			application.HandleEvent(event)
			event = application.GetNextEvent()
		}
	}
}

func (application *Application) WaitForEvent() tcell.Event {
	return <-application.screen.EventQ()
}

func (application *Application) GetNextEvent() tcell.Event {
	select {
	case event := <-application.screen.EventQ():
		return event

	default:
		return nil
	}
}

func (application *Application) HandleEvent(event tcell.Event) {
	switch event := event.(type) {
	case *tcell.EventResize:
		width, height := event.Size()

		slog.Debug(
			"screen resized",
			slog.Int("width", width),
			slog.Int("height", height),
		)

		application.size = tui.Size{
			Width:  width,
			Height: height,
		}

		message := tui.MsgResize{
			Size: application.size,
		}
		application.activeMode.Handle(message)

		application.screen.Sync()

	case *tcell.EventKey:
		translation := translateKey(event)
		slog.Debug("Key pressed", slog.String("key", translation))

		message := tui.MsgKey{
			Key: translateKey(event),
		}
		application.activeMode.Handle(message)

		// case *tcell.EventMouse:
		// 	x, y := event.Position()
		// 	position := tui.Position{X: x, Y: y}
		// 	clickHandler := grid.Get(position).OnClick

		// 	if clickHandler != nil && event.Buttons() == tcell.Button1 {
		// 		clickHandler()
		// 	}
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

// func (application *Application) createModel() {
// 	graph := application.graph

// 	// Create data sources
// 	input := data.NewVariable("")
// 	highlightedNodeIndex := data.NewVariable(0)
// 	selectedNodes := data.NewSliceList[*pkg.Node](nil)
// 	intersectionNodes := data.NewSliceList[*pkg.Node](nil)
// 	intersectionNodeNames := data.MapList(intersectionNodes, func(node *pkg.Node) string { return node.GetName() })

// 	// Cause intersection node list to be updated whenever the input or the selected nodes change
// 	updateIntersectionNodes := func() {
// 		nodes := determineIntersectionNodes(input.Get(), graph, data.CopyListToSlice(selectedNodes), true, true)
// 		intersectionNodes.SetSlice(nodes)
// 	}
// 	updateIntersectionNodes()
// 	data.DefineReaction(updateIntersectionNodes, input, selectedNodes)

// 	input.Observe(func() {
// 		if len(input.Get()) > 0 {
// 			application.updateHighlightedNode(input.Get())
// 		}
// 	})

// 	model := Model{
// 		input:                 input,
// 		selectedNodes:         selectedNodes,
// 		intersectionNodes:     intersectionNodes,
// 		highlightedNodeIndex:  highlightedNodeIndex,
// 		intersectionNodeNames: intersectionNodeNames,
// 	}

// 	application.model = model
// }

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

func (application *Application) updateHighlightedNode(target string) {
	intersectionNodes := application.model.IntersectionNodes().Get()
	nodes := list.ToSlice(intersectionNodes)

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

	update := application.model.Update()
	update.Highlight(bestMatchIndex)
	update.Apply()
}

func (application *Application) switchMode(mode mode) {
	application.activeMode = mode
	application.activeMode.Handle(tui.MsgResize{Size: application.size})
}

func (application *Application) selectHighlightedNode() {
	update := application.model.Update()
	update.SelectHighlightedNode()
	update.Apply()
}

func (application *Application) unselectLastNode() {
	update := application.model.Update()
	update.UnselectLastNode()
	update.Apply()
}

func (application *Application) highlight(index int) {
	update := application.model.Update()
	update.Highlight(index)
	update.Apply()
}

func (application *Application) selectHighlightedAndClearInput() {
	update := application.model.Update()
	update.SelectHighlightedNode()
	update.SetInput("")
	update.Apply()
}

func (application *Application) updateInput(newInput string) {
	update := application.model.Update()
	update.SetInput(strings.ToLower(newInput))
	update.Apply()
}
