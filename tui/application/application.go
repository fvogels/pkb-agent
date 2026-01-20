package application

import (
	"fmt"
	"log/slog"
	"os"
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/pkg/loaders/sequence"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/holder"
	"pkb-agent/tui/data"
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
	verbose          bool
	running          bool
	logFile          *os.File
	screen           tcell.Screen
	messageQueue     tui.MessageQueue
	screenSize       tui.Size
	graph            *pkg.Graph
	mode             mode
	model            data.Variable[*model.Model]
	activeModeHolder *holder.Component
	bindings         keyBindings
}

type keyBindings struct {
	mode data.Variable[list.List[tui.KeyBinding]]
	node data.Variable[list.List[tui.KeyBinding]]
	all  data.Value[list.List[tui.KeyBinding]]
}

type mode struct {
	view   *viewMode
	input  *inputMode
	active data.Variable[tui.Component]
}

func NewApplication(verbose bool) *Application {
	application := Application{
		verbose: verbose,
		running: true,
		logFile: nil,
	}

	createKeyBindings(&application.bindings)

	return &application
}

func createKeyBindings(bindings *keyBindings) {
	bindings.mode = data.NewVariable(list.New[tui.KeyBinding]())
	bindings.node = data.NewVariable(list.New[tui.KeyBinding]())

	bindings.all = data.MapValue2(
		&bindings.mode,
		&bindings.node,
		func(xs list.List[tui.KeyBinding], ys list.List[tui.KeyBinding]) list.List[tui.KeyBinding] {
			slog.Debug(
				"updating keybindings",
				"xs", list.String(xs, func(b tui.KeyBinding) string { return b.Key }),
				"ys", list.String(ys, func(b tui.KeyBinding) string { return b.Key }),
			)
			return list.Concatenate(xs, ys)
		},
	)
}

func (application *Application) Start() error {
	err := application.initializeLogging()
	if err != nil {
		return err
	}

	if err := application.initializeScreen(); err != nil {
		return err
	}

	application.createMessageQueue()

	if err := application.loadGraph(); err != nil {
		return err
	}

	application.model = data.NewVariable(model.New(application.graph))

	application.createModes(&application.mode)
	application.activeModeHolder = holder.New(application.messageQueue, &application.mode.active)

	application.mode.active.Set(application.mode.view)
	application.messageQueue.Enqueue(tui.MsgStateUpdated{})

	application.eventLoop()

	return nil
}

func (application *Application) createModes(mode *mode) {
	mode.view = newViewMode(application)
	mode.input = newInputMode(application)
	mode.active = data.NewVariable[tui.Component](nil)
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

		application.screenSize = tui.Size{
			Width:  width,
			Height: height,
		}

		message := tui.MsgResize{
			Size: application.screenSize,
		}

		application.handleMessage(message)

		application.screen.Sync()

	case *tcell.EventKey:
		translation := translateKey(event)
		slog.Debug("Key pressed", slog.String("key", translation))

		message := tui.MsgKey{
			Key: translateKey(event),
		}

		application.handleMessage(message)

	case *tui.EventMessage:
		application.handleMessage(event.Message)

		// case *tcell.EventMouse:
		// 	x, y := event.Position()
		// 	position := tui.Position{X: x, Y: y}
		// 	clickHandler := grid.Get(position).OnClick

		// 	if clickHandler != nil && event.Buttons() == tcell.Button1 {
		// 		clickHandler()
		// 	}
	}
}

func (application *Application) handleMessage(message tui.Message) {
	slog.Debug("Application handles message", slog.String("message", message.String()))

	switch message := message.(type) {
	case tui.MsgUpdateLayout:
		application.mode.active.Get().Handle(tui.MsgResize{
			Size: application.screenSize,
		})

	case messages.MsgQuit:
		application.running = false

	case messages.MsgSelectHighlightedNode:
		application.selectHighlightedNode()

	case messages.MsgUnselectLastNode:
		application.unselectLastNode()

	case messages.MsgSetModeKeyBindings:
		application.bindings.mode.Set(message.Bindings)

	case messages.MsgSetNodeKeyBindings:
		application.bindings.node.Set(message.Bindings)

	case messages.MsgActivateInputMode:
		application.switchMode(application.mode.input)

	case messages.MsgSwitchLinksView:
		application.switchLinksView()

	case messages.MsgLockSelectedNodes:
		application.lockSelectedNodes()

	case messages.MsgUnlockSelectedNodes:
		application.unlockSelectedNodes()

	case tui.MsgCommand:
		message.Command()

	default:
		application.mode.active.Get().Handle(message)
	}
}

func (application *Application) Render() {
	screen := application.screen
	activeMode := application.mode.active

	screen.Clear()

	grid := activeMode.Get().Render()
	gridSize := grid.Size()
	runes := make([]rune, 1)

	timeBeforeUpdate := time.Now()

	for y := range gridSize.Height {
		for x := range gridSize.Width {
			position := tui.Position{X: x, Y: y}
			cell := grid.At(position)
			runes[0] = cell.Contents
			screen.Put(x, y, string(runes), *cell.Style)
		}
	}

	screen.Show()
	slog.Debug("Screen updated", slog.String("duration", time.Since(timeBeforeUpdate).String()))
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

func (application *Application) loadGraph() error {
	before := time.Now()
	loader := sequence.New()
	path := pathlib.New(`F:\repos\pkb\pkb-data\root.yaml`)

	graph, err := pkg.LoadGraph(path, loader)
	if err != nil {
		slog.Debug("Error occurred", slog.String("message", err.Error()))
		return err
	}
	slog.Debug("Graph loaded", slog.String("loadTime", time.Since(before).String()))

	application.graph = graph

	return nil
}

func (application *Application) findIndexOfIntersectionNode(intersectionNodes []*pkg.Node, target string) int {
	bestMatchIndex, found := slices.BinarySearchFunc(
		intersectionNodes,
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

	return bestMatchIndex
}

func (application *Application) switchMode(mode tui.Component) {
	application.mode.active.Set(mode)
	application.messageQueue.Enqueue(tui.MsgResize{Size: application.screenSize})
}

func (application *Application) updateModel(updater func(model *model.Model)) {
	updatedModel := *application.model.Get()
	updater(&updatedModel)
	application.model.Set(&updatedModel)
	application.messageQueue.Enqueue(tui.MsgStateUpdated{})
}

func (application *Application) switchLinksView() {
	application.updateModel(func(model *model.Model) {
		model.ShowNodeLinks = !model.ShowNodeLinks
	})
}

func (application *Application) selectHighlightedNode() {
	application.updateModel(func(model *model.Model) {
		model.SelectHighlightedNode()
		model.DetermineIntersectionNodes()
		model.HighlightedNodeIndex = 0
	})
}

func (application *Application) unselectLastNode() {
	application.updateModel(func(model *model.Model) {
		model.UnselectLastNode()
		model.DetermineIntersectionNodes()
	})
}

func (application *Application) highlight(index int) {
	application.updateModel(func(model *model.Model) {
		model.HighlightedNodeIndex = index
		model.DetermineIntersectionNodes()
	})
}

func (application *Application) selectHighlightedAndClearInput() {
	application.updateModel(func(model *model.Model) {
		if model.IntersectionNodes.Size() > 0 {
			model.SelectHighlightedNode()
		}
		model.Input = ""
		model.DetermineIntersectionNodes()
	})
}

func (application *Application) updateInputAndHighlightBestMatch(newInput string) {
	lowerCasedNewInput := strings.ToLower(newInput)

	application.updateModel(func(model *model.Model) {
		model.Input = newInput
		model.DetermineIntersectionNodes()
		model.HighlightedNodeIndex = application.findIndexOfIntersectionNode(list.ToSlice(model.IntersectionNodes), lowerCasedNewInput)
	})
}

func (application *Application) createMessageQueue() {
	application.messageQueue = tui.NewMessageQueue(application.screen.EventQ())
}

func (application *Application) lockSelectedNodes() {
	application.updateModel(func(model *model.Model) {
		model.LockSelectedNodes()
	})
}

func (application *Application) unlockSelectedNodes() {
	application.updateModel(func(model *model.Model) {
		model.UnlockSelectedNodes()
	})
}
