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
	"reflect"
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
	size             tui.Size
	graph            *pkg.Graph
	model            model.Model
	viewMode         *viewMode
	inputMode        *inputMode
	activeMode       data.Variable[tui.Component]
	activeModeHolder *holder.Component
	modeBindings     data.Variable[list.List[tui.KeyBinding]]
	nodeBindings     data.Variable[list.List[tui.KeyBinding]]
	keyBindings      data.Value[list.List[tui.KeyBinding]]
}

func NewApplication(verbose bool) *Application {
	application := Application{
		verbose:      verbose,
		running:      true,
		logFile:      nil,
		activeMode:   data.NewVariable[tui.Component](nil),
		modeBindings: data.NewVariable(list.New[tui.KeyBinding]()),
		nodeBindings: data.NewVariable(list.New[tui.KeyBinding]()),
	}

	application.keyBindings = data.MapValue2(
		&application.modeBindings,
		&application.nodeBindings,
		func(xs list.List[tui.KeyBinding], ys list.List[tui.KeyBinding]) list.List[tui.KeyBinding] {
			return list.Concatenate(xs, ys)
		},
	)

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

	application.createMessageQueue()

	if err := application.loadGraph(); err != nil {
		return err
	}

	application.model = model.New(application.graph)
	application.viewMode = newViewMode(application)
	application.inputMode = newInputMode(application)
	application.activeModeHolder = holder.New(application.messageQueue, &application.activeMode)

	application.switchMode(application.viewMode)

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
		slog.Debug("Application handles message", slog.String("messageType", reflect.TypeOf(message).String()))

		application.activeMode.Get().Handle(message)
		application.screen.Sync()

	case *tcell.EventKey:
		translation := translateKey(event)
		slog.Debug("Key pressed", slog.String("key", translation))

		message := tui.MsgKey{
			Key: translateKey(event),
		}

		slog.Debug("Application handles message", slog.String("messageType", reflect.TypeOf(message).String()))

		application.activeMode.Get().Handle(message)

	case *tui.EventMessage:
		message := event.Message

		slog.Debug("Application handles message", slog.String("messageType", reflect.TypeOf(message).String()))

		switch message := message.(type) {
		case tui.MsgUpdateLayout:
			application.activeMode.Get().Handle(tui.MsgResize{
				Size: application.size,
			})

		case messages.MsgQuit:
			application.running = false

		case messages.MsgSelectHighlightedNode:
			application.selectHighlightedNode()

		case messages.MsgUnselectLastNode:
			application.unselectLastNode()

		case messages.MsgSetModeKeyBindings:
			application.modeBindings.Set(message.Bindings)

		case messages.MsgSetNodeKeyBindings:
			application.nodeBindings.Set(message.Bindings)

		case messages.MsgActivateInputMode:
			application.switchMode(application.inputMode)

		default:
			application.activeMode.Get().Handle(message)
		}

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

	grid := activeMode.Get().Render()
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
	application.activeMode.Set(mode)
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

func (application *Application) updateInputAndHighlightBestMatch(newInput string) {
	lowerCasedNewInput := strings.ToLower(newInput)

	update := application.model.Update()
	update.SetInput(lowerCasedNewInput)
	update.DetermineIntersectionNodes()
	intersectionNodes := update.UpdatedModel.IntersectionNodes().Get()
	update.Highlight(application.findIndexOfIntersectionNode(list.ToSlice(intersectionNodes), lowerCasedNewInput))
	update.Apply()
}

func (application *Application) createMessageQueue() {
	application.messageQueue = tui.NewMessageQueue(application.screen.EventQ())
}
