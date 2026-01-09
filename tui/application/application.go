package application

import (
	"fmt"
	"log/slog"
	"os"
	"pkb-agent/graph/loaders/sequence"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/component/stringlist"
	"pkb-agent/tui/component/stringsview"
	"pkb-agent/tui/data"
	"pkb-agent/util/pathlib"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Application struct {
	verbose bool
	logFile *os.File
	screen  tcell.Screen
	graph   *pkg.Graph
	model   Model
}

type Model struct {
	input                 *data.Variable[string]
	selectedNodes         *data.SliceList[*pkg.Node]
	intersectionNodes     *data.SliceList[*pkg.Node]
	selectedItemIndex     *data.Variable[int]
	intersectionNodeNames data.List[string]
}

type View struct {
	intersectionNodeView *stringlist.Component
	root                 tui.Component
}

func NewApplication(verbose bool) *Application {
	application := Application{
		verbose: verbose,
	}

	return &application
}

func (application *Application) Start() error {
	// out, _ := os.Create("profile.txt")
	// pprof.StartCPUProfile(out)
	// defer pprof.StopCPUProfile()

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

	application.model = application.createModel()

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

func createViews(model *Model) *View {
	intersectionNodeView := stringlist.New(model.intersectionNodeNames, model.selectedItemIndex)
	intersectionNodeView.SetOnSelectionChanged(func(value int) { model.selectedItemIndex.Set(value) })
	root := intersectionNodeView

	view := View{
		intersectionNodeView: intersectionNodeView,
		root:                 root,
	}

	return &view
}

func (application *Application) eventLoop() {
	// style := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	screen := application.screen
	model := &application.model

	// Views
	intersectionNodeView := stringlist.New(model.intersectionNodeNames, model.selectedItemIndex)
	intersectionNodeView.SetOnSelectionChanged(func(value int) { model.selectedItemIndex.Set(value) })

	root := intersectionNodeView

	for {
		// Update screen
		screen.Clear()

		grid := root.Render()
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

		// Poll event (this can be in a select statement as well)
		ev := <-screen.EventQ()

		// Process event
		switch event := ev.(type) {
		case *tcell.EventResize:
			width, height := event.Size()

			root.Handle(tui.MsgResize{
				Size: tui.Size{
					Width:  width,
					Height: height,
				},
			})
			screen.Sync()

		case *tcell.EventKey:
			if event.Str() == "q" {
				return
			} else {
				translation := translateKey(event)

				slog.Debug("Key pressed", slog.String("key", translation))

				switch translation {
				default:
					message := tui.MsgKey{
						Key: translateKey(event),
					}

					root.Handle(message)
				}
			}

		case *tcell.EventMouse:
			x, y := event.Position()
			position := tui.Position{X: x, Y: y}
			clickHandler := grid.Get(position).OnClick

			if clickHandler != nil && event.Buttons() == tcell.Button1 {
				clickHandler()
			}

			// switch ev.Buttons() {
			// case tcell.Button1, tcell.Button2:

			// case tcell.ButtonNone:

			// }
		}
	}
}

func (application *Application) createModel() Model {
	graph := application.graph

	input := data.NewVariable("")
	selectedNodes := data.NewSliceList[*pkg.Node](nil)
	intersectionNodes := data.NewSliceList[*pkg.Node](nil)
	updateIntersectionNodes := func() {
		nodes := determineIntersectionNodes(input.Get(), graph, data.CopyListToSlice(selectedNodes), true, true)
		intersectionNodes.SetSlice(nodes)
	}
	updateIntersectionNodes()
	data.DefineReaction(updateIntersectionNodes, input, selectedNodes)
	selectedItemIndex := data.NewVariable(0)
	intersectionNodeNames := data.MapList(intersectionNodes, func(node *pkg.Node) string { return node.GetName() })

	return Model{
		input:                 input,
		selectedNodes:         selectedNodes,
		intersectionNodes:     intersectionNodes,
		selectedItemIndex:     selectedItemIndex,
		intersectionNodeNames: intersectionNodeNames,
	}
}

func (application *Application) initializeLogging() error {
	if application.verbose {
		logFile, err := os.Create("ui.log")
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
