package main

import (
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/nelsonmarro/fine-markdown-editor/internal/theme"
)

type Config struct {
	EditWidget   *widget.Entry
	PrevWidget   *widget.RichText
	CurrentFile  fyne.URI
	SaveMenuItem *fyne.MenuItem
}

var cfg Config

func main() {
	// create a fyne app
	a := app.New()

	a.Settings().SetTheme(&theme.MyTheme{})

	// create a new window for the app
	w := a.NewWindow("Markdown Editor")

	// get the user interface
	edit, preview := cfg.makeUI()
	cfg.createMenuItems(w)

	// set the content of the window
	w.SetContent(container.NewHSplit(edit, preview))

	// show window and run the app
	w.Resize(fyne.Size{Width: 800, Height: 500})
	w.CenterOnScreen()
	w.ShowAndRun()
}

func (app *Config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")
	app.EditWidget = edit
	app.PrevWidget = preview

	// Add an event listener OnChanged in the edit widget
	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (app Config) createMenuItems(w fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open...", app.openFunc(w))
	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(w))
	saveAsMenuItem := fyne.NewMenuItem("Save As...", app.saveAsFunc(w))

	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	menu := fyne.NewMainMenu(fileMenu)

	w.SetMainMenu(menu)
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *Config) saveFunc(w fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			writer, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			writer.Write([]byte(app.EditWidget.Text))
			defer writer.Close()
		}
	}
}

func (app *Config) openFunc(w fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// user click cancel
			if reader == nil {
				return
			}

			defer reader.Close()

			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, w)
			}

			app.EditWidget.SetText(string(data))

			app.CurrentFile = reader.URI()
			w.SetTitle(w.Title() + " - " + reader.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, w)

		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (app *Config) saveAsFunc(w fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			if write == nil {
				// the user cancelled the dialog
				return
			}

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Please name yout file with a .md extension!", w)
				return
			}

			// save file
			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()

			defer write.Close()

			w.SetTitle(w.Title() + " - " + write.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, w)

		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)

		saveDialog.Show()
	}
}
