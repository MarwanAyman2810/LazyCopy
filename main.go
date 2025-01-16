package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
)

func getSelectedText() string {
	cmd := exec.Command("xclip", "-o", "-selection", "primary")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting selection: %v\n", err)
		return ""
	}
	return strings.TrimSpace(string(output))
}

func main() {

	myApp := app.New()

	// Create the copy window
	copyWindow := myApp.NewWindow("Copy/Paste")
	copyWindow.Resize(fyne.NewSize(50, 50))
	copyWindow.SetFixedSize(true)
	copyWindow.Canvas().Focused()
	copyWindow.Hide()

	var selectedText string
	var lastSelection string

	// Add label to display selected text
	textLabel := widget.NewLabel("")
	copyBtn := widget.NewButton("Copy", func() {
		if selectedText != "" {
			if err := clipboard.WriteAll(selectedText); err != nil {
				fmt.Printf("Error writing to clipboard: %v\n", err)
			}
			// Window stays open and on top
		}
	})
	copyBtn.Resize(fyne.NewSize(100, copyBtn.MinSize().Height))

	// Add paste button
	pasteBtn := widget.NewButton("Paste", func() {
		copyWindow.Hide()
		// Give a small delay for window focus to change
		time.Sleep(100 * time.Millisecond)
		robotgo.KeyTap("v", "ctrl")
	})
	pasteBtn.Resize(fyne.NewSize(100, pasteBtn.MinSize().Height))

	// Add paste and go button
	pasteGoBtn := widget.NewButton("Paste & Go", func() {
		copyWindow.Hide()
		// Give a small delay for window focus to change
		time.Sleep(100 * time.Millisecond)
		robotgo.KeyTap("v", "ctrl")
		time.Sleep(50 * time.Millisecond)
		robotgo.KeyTap("enter")
	})
	pasteGoBtn.Resize(fyne.NewSize(100, pasteGoBtn.MinSize().Height))

	// Create a 1x3 grid layout for buttons (1 row, 3 columns)
	buttonsBox := container.New(layout.NewGridLayout(3), copyBtn, pasteBtn, pasteGoBtn)

	content := container.NewVBox(
		textLabel,
		buttonsBox,
	)
	copyWindow.SetContent(container.NewPadded(content))

	// Monitor text selection
	go func() {
		// Initial state
		lastSelection = getSelectedText()
		isWindowVisible := false

		for {
			newSelection := getSelectedText()
			// Only show window if there's a new selection different from the last one
			if newSelection != "" && newSelection != lastSelection {
				selectedText = newSelection
				lastSelection = newSelection
				textLabel.SetText(selectedText)
				copyWindow.Show()
				copyWindow.CenterOnScreen()
				copyWindow.RequestFocus()
				isWindowVisible = true
			} else if newSelection == "" && isWindowVisible {
				copyWindow.Hide()
				isWindowVisible = false
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	myApp.Run()
}
