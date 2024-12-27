package main

import (
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

func Test_makeUI(t *testing.T) {
	var testCfg Config

	edit, preview := testCfg.makeUI()

	test.Type(edit, "Hello")

	if preview.String() != "Hello" {
		t.Error("Failed -- preview did not update with the correct value")
	}
}

func Test_RunApp(t *testing.T) {
	var testCfg Config

	testApp := test.NewApp()
	testWin := testApp.NewWindow("Test Markdown")

	edit, preview := testCfg.makeUI()

	testCfg.createMenuItems(testWin)

	testWin.SetContent(container.NewHSplit(edit, preview))

	testApp.Run()

	test.Type(edit, "Some Text")
	if preview.String() != "Some Text" {
		t.Error("Failed -- preview did not update with the correct value")
	}
}
