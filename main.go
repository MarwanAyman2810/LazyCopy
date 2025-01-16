package main

import (
	"fmt"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
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
	copyWindow := myApp.NewWindow("Copy")
	copyWindow.Resize(fyne.NewSize(300, 100))
	copyWindow.SetFixedSize(true)
	copyWindow.Hide()

	var selectedText string

	// Add label to display selected text
	textLabel := widget.NewLabel("")
	copyBtn := widget.NewButton("Copy", func() {
		if selectedText != "" {
			if err := clipboard.WriteAll(selectedText); err != nil {
				fmt.Printf("Error writing to clipboard: %v\n", err)
			}
		}
		copyWindow.Hide()
	})

	// Use vertical box to stack label and button
	content := container.NewVBox(
		textLabel,
		copyBtn,
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
		}
	}()

	myApp.Run()
}
