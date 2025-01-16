package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
	copyWindow.Resize(fyne.NewSize(300, 150))
	copyWindow.SetFixedSize(true)
	copyWindow.Canvas().Focused()
	copyWindow.Hide()

	var selectedText string

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

	// Add paste button
	pasteBtn := widget.NewButton("Paste", func() {
		copyWindow.Hide()
		// Give a small delay for window focus to change
		time.Sleep(100 * time.Millisecond)
		robotgo.KeyTap("v", "ctrl")
	})

	// Use vertical box to stack label and buttons
	content := container.NewVBox(
		textLabel,
		copyBtn,
		pasteBtn,
	)
	copyWindow.SetContent(container.NewPadded(content))

	// Monitor text selection
	go func() {
		for {
			newSelection := getSelectedText()
			if newSelection != "" && newSelection != selectedText {
				selectedText = newSelection
				textLabel.SetText(selectedText)
				copyWindow.Show()
				copyWindow.CenterOnScreen()
				copyWindow.RequestFocus()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	myApp.Run()
}
