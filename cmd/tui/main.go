package main

import (
	"gps/internal/config"
	"gps/internal/gpt"
	"path/filepath"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		panic(err)
	}

	app := tview.NewApplication()

	// the prompt area has a text area and a help hint
	promptTextArea := createPromptTextArea()
	promptView := createPromptView(promptTextArea)
	previewTree := createPreviewTree()
	previewTextView := tview.NewTextView()

	pages := tview.NewPages()
	helpPage := createHelpPage(pages)
	waitPage := tview.NewModal().
		SetText("Processing... Please wait...")
		// AddButtons([]string{"Abort"}).
		// SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		// 	app.Stop()
		// })

	// Frames in a Flex layout.
	mainPage := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(promptView, 0, 1, true).
		AddItem(tview.NewFlex().
			AddItem(previewTree, 0, 1, false).
			AddItem(previewTextView, 0, 3, false), 0, 7, false)

	pages.AddAndSwitchToPage("main", mainPage, true).
		AddPage("help", helpPage, true, false).
		AddPage("wait", waitPage, false, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlH {
			pages.ShowPage("help") //TODO: Check when clicking outside help window with the mouse. Then clicking help again.
			return nil
		} else if event.Key() == tcell.KeyEnter {
			pages.ShowPage("wait")
			go func() {
				app.QueueUpdateDraw(func() {
					opts := gpt.ProjectOptions{
						Name:   "test",
						Prompt: promptTextArea.GetText(),
					}
					resp, err := gpt.NewProject(cfg.Gpt, &opts)
					previewTextView.SetText(resp.Preview)
					if err != nil {
						panic(err)
					}
					sort.Strings(resp.Files)
					// fmt.Println(resp.Files)
					refreshPreviewTree(previewTree, resp.Files)
					pages.HidePage("wait")
				})
			}()
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func createPromptTextArea() *tview.TextArea {
	textArea := tview.NewTextArea().
		SetPlaceholder("Describe your application here. eg: 'The application is called voyager and contains two services: janeway and chakotay'.")
	textArea.SetTitle("Prompt").
		SetBorder(true)

	return textArea
}

// Create a tview box with the title "Prompt"
func createPromptView(textArea *tview.TextArea) *tview.Flex {
	helpInfo := tview.NewTextView().
		SetText(" Enter: Submit prompt | Ctrl-H: Show Help | Ctrl-C: Exit")

	promptView := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textArea, 0, 6, true).
		AddItem(helpInfo, 0, 1, false)

	return promptView
}

func addNodes(path string, nodes *map[string]*tview.TreeNode) {
	dir := filepath.Dir(path)
	if _, exists := (*nodes)[dir]; !exists {
		addNodes(dir, nodes)
	}
	(*nodes)[path] = tview.NewTreeNode(filepath.Base(path)).SetReference(path)
	(*nodes)[dir].AddChild((*nodes)[path])
}

func refreshPreviewTree(tree *tview.TreeView, filepaths []string) *map[string]*tview.TreeNode {
	nodes := make(map[string]*tview.TreeNode)
	nodes["."] = tview.NewTreeNode(".").SetReference(".").SetSelectable(false)

	for _, path := range filepaths {
		addNodes(path, &nodes)
	}

	tree.SetRoot(nodes["."]).SetCurrentNode(nodes["."])

	return &nodes
}

// Show a navigable tree view of the current directory.
func createPreviewTree() *tview.TreeView {
	tree := tview.NewTreeView()

	refreshPreviewTree(tree, []string{})

	return tree
}

const editorHelp = `[green]Navigation

	[yellow]Left arrow[white]: Move left.
	[yellow]Right arrow[white]: Move right.
	[yellow]Down arrow[white]: Move down.
	[yellow]Up arrow[white]: Move up.
	[yellow]Ctrl-A, Home[white]: Move to the beginning of the current line.
	[yellow]Ctrl-E, End[white]: Move to the end of the current line.
	[yellow]Ctrl-F, page down[white]: Move down by one page.
	[yellow]Ctrl-B, page up[white]: Move up by one page.
	[yellow]Alt-Up arrow[white]: Scroll the page up.
	[yellow]Alt-Down arrow[white]: Scroll the page down.
	[yellow]Alt-Left arrow[white]: Scroll the page to the left.
	[yellow]Alt-Right arrow[white]:  Scroll the page to the right.
	[yellow]Alt-B, Ctrl-Left arrow[white]: Move back by one word.
	[yellow]Alt-F, Ctrl-Right arrow[white]: Move forward by one word.

[green]Editing[white]

	Type to enter text.
	[yellow]Ctrl-K[white]: Delete until the end of the line.
	[yellow]Ctrl-W[white]: Delete the rest of the word.
	[yellow]Ctrl-U[white]: Delete the current line.

[green]Selecting Text[white]

	Move while holding Shift or drag the mouse.
	Double-click to select a word.

[green]Clipboard

	[yellow]Ctrl-Q[white]: Copy.
	[yellow]Ctrl-X[white]: Cut.
	[yellow]Ctrl-V[white]: Paste.
	
[green]Undo

	[yellow]Ctrl-Z[white]: Undo.
	[yellow]Ctrl-Y[white]: Redo.

[blue]Press Escape or Enter to return.`

func createHelpPage(pages *tview.Pages) *tview.Frame {
	helpText := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetText(editorHelp)

	help := tview.NewFrame(helpText).
		SetBorders(1, 1, 0, 0, 2, 2)

	help.SetBorder(true).
		SetTitle("Help").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter {
				pages.SwitchToPage("main")
				return nil
			}
			return event
		})

	return help
}
