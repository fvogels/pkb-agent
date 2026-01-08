package tuimain

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"pkb-agent/tui"
	"pkb-agent/tui/component/border"
	"pkb-agent/tui/component/cache"
	"pkb-agent/tui/component/docksouth"
	"pkb-agent/tui/component/input"
	"pkb-agent/tui/component/stringlist"
	"pkb-agent/tui/component/stringsview"
	"pkb-agent/tui/data"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func Start(verbose bool) error {
	// out, _ := os.Create("profile.txt")
	// pprof.StartCPUProfile(out)
	// defer pprof.StopCPUProfile()

	closeLog, err := initializeLogging(verbose)
	if err != nil {
		return err
	}
	defer closeLog()

	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	// boxStyle := tcell.StyleDefault.Foreground(color.White).Background(color.Purple)

	// Initialize screen
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	screen.SetStyle(defStyle)
	screen.EnableMouse()
	// s.EnablePaste()
	screen.Clear()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Event loop
	eventLoop(screen)

	return nil
}

func eventLoop(screen tcell.Screen) {
	style := tcell.StyleDefault.Background(color.Green).Foreground(color.Reset)
	statusStyle := tcell.StyleDefault.Background(color.Red).Foreground(color.Reset)

	strings := []string{"a", "bb", "ccc", "dddd", "eeeee", "ffffff", "ggggggg", "hhhhhhhh", "a", "bb", "ccc", "dddd", "eeeee", "ffffff", "ggggggg", "hhhhhhhh", "a", "bb", "ccc", "dddd", "eeeee", "ffffff", "ggggggg", "hhhhhhhh", "a", "bb", "ccc", "dddd", "eeeee", "ffffff", "ggggggg", "hhhhhhhh"}
	// runes := util.Map(strings, func(s string) []rune { return []rune(s) })
	// items := data.NewSliceList(strings)
	// runedItems := ItemList{
	// 	items: runes,
	// 	style: &style,
	// }
	// selectedItem := data.NewVariable(0)

	text := data.NewVariable("")
	selectedItem := data.NewVariable(0)

	list := stringlist.New(data.NewSliceList(strings), selectedItem)
	list.SetOnSelectionChanged(func(index int) { selectedItem.Set(index) })

	mainView := border.New(cache.New(list, selectedItem), style)

	statusBar := input.New(text, statusStyle, func(s string) { text.Set(s) })
	root := docksouth.New(mainView, statusBar, 1)

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
		slog.Debug("stringlist", slog.String("duration", time.Since(timeBeforeUpdate).String()))

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

func initializeLogging(verbose bool) (func(), error) {
	if verbose {
		logFile, err := os.Create("ui.log")
		if err != nil {
			fmt.Println("Failed to create log")
			return nil, err
		}

		logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)

		return func() { logFile.Close() }, nil
	}

	return func() {}, nil
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
