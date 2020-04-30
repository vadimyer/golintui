package gui

import (
	"github.com/rivo/tview"

	"github.com/nakabonne/golintui/pkg/gui/item"
)

type Gui struct {
	application *tview.Application
	pages       *tview.Pages

	lintersItem     *item.Linters
	sourceFilesItem *item.SourceFiles
	resultsItem     *item.Results
	infoItem        *item.Info
}

func New() *Gui {
	return &Gui{
		application:     tview.NewApplication(),
		lintersItem:     item.NewLinters(),
		sourceFilesItem: item.NewSourceFiles("."),
		resultsItem:     item.NewResults(),
		infoItem:        item.NewInfo(),
	}
}

func (g *Gui) Run() error {
	g.setKeybind()
	g.initGrid()
	if err := g.application.Run(); err != nil {
		g.application.Stop()
		return err
	}
	return nil
}

// initGrid sets a grid based layout as a root primitive for the application.
func (g *Gui) initGrid() {
	grid := tview.NewGrid().
		SetRows(1, 0).
		SetColumns(30, 40, 0, 0).
		SetBorders(true).
		AddItem(g.infoItem, 0, 0, 1, 2, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(g.lintersItem, 1, 0, 1, 1, 0, 100, false).
		AddItem(g.sourceFilesItem, 1, 1, 1, 1, 0, 100, true).
		AddItem(g.resultsItem, 0, 2, 2, 2, 0, 100, false)

	g.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)
	g.application.SetRoot(g.pages, true)
}