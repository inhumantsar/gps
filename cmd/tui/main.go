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

	pages := tview.NewPages()
	helpPage := createHelpPage(pages)
	waitPage := tview.NewModal().
		SetText("Processing... Please wait...")
		// AddButtons([]string{"Abort"}).
		// SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		// 	app.Stop()
		// })

	previewTree := createPreviewTree()

	initForm := createInitializerView(app,
		cfg,
		func() { pages.ShowPage("wait") },
		func(resp *gpt.InitializeResponse) {
			sort.Strings(resp.Files)
			// fmt.Println(resp.Files)
			refreshPreviewTree(previewTree, resp.Files)
			pages.HidePage("wait")
		})

	// Frames in a Flex layout.
	mainPage := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(initForm, 0, 1, true).
		AddItem(previewTree, 0, 8, false)
		// AddItem(previewTextView, 0, 3, false), 0, 7, false)

	previewTree.SetBorder(true).SetTitle("Preview")

	pages.AddAndSwitchToPage("main", mainPage, true).
		AddPage("help", helpPage, true, false).
		AddPage("wait", waitPage, false, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlH {
			pages.ShowPage("help") //TODO: Check when clicking outside help window with the mouse. Then clicking help again.
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func createInitializerView(app *tview.Application, cfg *config.Config, pre func(), post func(*gpt.InitializeResponse)) *tview.Grid {
	p := "Wafflemart allows users to order waffles. It consists of an API ('batter') and a worker ('iron'). batter runs on AWS Lambda and the worker runs on ECS."
	promptTextArea := tview.NewTextArea().
		SetPlaceholder(p).
		SetWordWrap(true).
		SetSize(4, 0)
		// SetTextStyle(tcell.StyleDefault.
		// 	Foreground(tcell.ColorDarkSlateGray)).
		// 	Foreground(tcell.ColorDimGray)).
		// SetPlaceholderStyle(tcell.StyleDefault.
		// 	Foreground(tcell.ColorDarkSlateGray))

	promptTextArea.SetBorder(true).
		SetTitle(" Describe your application  ").
		// SetBackgroundColor(tcell.ColorDarkSlateGray).
		// SetTitleColor(tcell.ColorDarkSlateGray).
		SetTitleAlign(tview.AlignLeft)

	// promptLabel := tview.NewTextView().SetText("Describe your application:")

	save := func() {
		pre()
		go func() {
			app.QueueUpdateDraw(func() {
				opts := gpt.InitializeOptions{
					Stream: false,
					Prompt: promptTextArea.GetText(),
				}
				resp, err := gpt.Initialize(cfg.Gpt, &opts)
				// previewTextView.SetText(resp.Preview)
				if err != nil {
					panic(err)
				}
				post(resp)
			})
		}()
	}
	saveBtn := tview.NewButton("Build").
		SetSelectedFunc(save).
		SetStyle(tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorDimGray))
	// saveBtn.
	// 	SetBorder(true)

	quitBtn := tview.NewButton("Quit").
		SetSelectedFunc(func() {
			app.Stop()
		}).
		SetStyle(tcell.StyleDefault.Background(tcell.ColorDarkRed).Foreground(tcell.ColorDimGray))
	// quitBtn.
	// 	SetBorder(true)

	form := tview.NewGrid().
		// SetDirection(tview.FlexRow).
		// AddItem(promptLabel, 0, 1, true).
		// Setcol
		SetColumns(-2, 1, 15).
		AddItem(promptTextArea, 0, 0, 1, 1, 0, 0, true).
		AddItem(tview.NewBox(), 0, 1, 1, 1, 0, 0, false).
		AddItem(tview.NewGrid().
			// SetDirection(tview.FlexRow).
			// SetColumns(10, 2, 10, 2, 10, 2, 10, 2, 10, 2, 10).
			// SetMinSize(1, 10).
			// SetColumns(0, 10, 0, 10, 0, 10, 0, 10, 0, 10, 0).
			// AddItem(tview.NewBox(), 0, 10, false).
			SetRows(0, -3, 0, -3, 0).
			AddItem(saveBtn, 1, 0, 1, 1, 0, 0, false).
			// AddItem(tview.NewBox(), 0, 1, false).
			AddItem(quitBtn, 3, 0, 1, 1, 0, 0, false), 0, 2, 1, 1, 0, 0, false)

	// form.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
	// form.SetLabelColor(tcell.ColorDimGray)

	// helpText := tview.NewTextView().
	// 	SetText("These prompts will be used to generate (or [::u]re[::-]generate) a multi-service application.").
	// 	SetDynamicColors(true)

	// view := tview.NewFlex().
	// 	SetDirection(tview.FlexRow).
	// 	AddItem(tview.NewBox(), 0, 1, false).
	// 	// AddItem(helpText, 3, 0, false).
	// 	AddItem(form, 0, 1, false)

	return form
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
