package item

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

const (
	DefaultDirColor  = tcell.ColorSilver
	SelectedDirColor = tcell.ColorLime
)

type SourceFiles struct {
	*tview.TreeView
	logger *logrus.Entry
}

func NewSourceFiles(logger *logrus.Entry, rootDir string) *SourceFiles {
	s := &SourceFiles{
		TreeView: tview.NewTreeView(),
		logger:   logger,
	}

	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorWhite)
	s.SetRoot(root).
		SetCurrentNode(root).
		SetBorder(true).
		SetTitle("Source Files").
		SetTitleAlign(tview.AlignLeft)

	if err := s.addChildren(root, rootDir); err != nil {
		panic(err) // TODO: Emit log instead of panic
	}
	return s
}

func (s *SourceFiles) SetKeybinds(globalKeybind func(event *tcell.EventKey), selectAction, unselectAction func(*tview.TreeNode, string)) {
	s.SetSelectedFunc(func(node *tview.TreeNode) {
		ref, ok := node.GetReference().(string)
		if !ok {
			return
		}
		if node.GetColor() == SelectedDirColor {
			unselectAction(node, ref)
		} else {
			selectAction(node, ref)
		}
	})

	s.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		node := s.GetCurrentNode()
		switch event.Rune() {
		case 'o':
			s.SwitchToggle(node)
		}
		globalKeybind(event)
		return event
	})
}

// AddChildren adds child nodes to the given node which represents a directory.
func (s *SourceFiles) addChildren(node *tview.TreeNode, path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		child := tview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(path, file.Name())).
			SetSelectable(file.IsDir()).SetColor(DefaultDirColor)
		node.AddChild(child)
	}
	return nil
}

// SwitchToggle switches the current display state.
func (s *SourceFiles) SwitchToggle(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}
	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		path := reference.(string)
		if err := s.addChildren(node, path); err != nil {
			s.logger.Error(err)
		}
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}
