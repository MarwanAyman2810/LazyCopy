// main.go
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Entry point of the application
func main() {
	// Create a new Fyne application
	myApp := app.New()

	// Create a new window and set its title
	myWindow := myApp.NewWindow("LazyCopy")

	// Create a new multiline entry widget
	textEntry := widget.NewMultiLineEntry()
	textEntry.SetPlaceHolder("Select some text here...")

	// Create a label to show debug messages
	debugLabel := widget.NewLabel("Debug messages will appear here")

	// Add a listener to detect selection changes
	textEntry.OnCursorChanged = func() {
		selectedText := textEntry.SelectedText()
		if selectedText != "" {
			debugLabel.SetText("Selected: " + selectedText)
		} else {
			debugLabel.SetText("No text selected")
		}
	}

	// Make the debug label more visible
	debugLabel.Alignment = fyne.TextAlignCenter
	debugLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Set the content with some padding
	content := container.NewVBox(
		textEntry,
		widget.NewSeparator(), // Add a visual separator
		debugLabel,
	)

	// Set the content of the window to include both widgets
	myWindow.SetContent(content)

	// Show the window and start the application
	myWindow.ShowAndRun()
}
