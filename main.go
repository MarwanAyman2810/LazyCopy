package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
	hook "github.com/robotn/gohook"
)

func main() {
	myApp := app.New()

	// Create only the copy button window
	copyWindow := myApp.NewWindow("Copy")
	copyWindow.Resize(fyne.NewSize(100, 40))
	copyWindow.SetFixedSize(true)
	copyWindow.Hide()

	copyBtn := widget.NewButton("Copy", func() {
		if text, err := clipboard.ReadAll(); err == nil {
			fmt.Printf("Copied: %q\n", text)
		}
		copyWindow.Hide()
	})
	copyWindow.SetContent(container.NewPadded(copyBtn))

	// Create a hidden main window (required by Fyne)
	mainWindow := myApp.NewWindow("Main")
	mainWindow.Hide()

	// Initialize the hook
	hooks := hook.Start() // Changed from gohook.NewHook() to hook.Start()
	defer hook.End()      // Changed from hooks.Unhook() to hook.End()

	lastClick := time.Now()

	// Register mouse down event
	hook.Register(hook.MouseDown, []string{}, func(e hook.Event) {
		now := time.Now()
		if now.Sub(lastClick) < 500*time.Millisecond {
			// Show the copy window on double-click
			copyWindow.CenterOnScreen()
			copyWindow.Show()
			copyWindow.RequestFocus()
		}
		lastClick = now
	})

	// Start processing events in a separate goroutine
	go func() {
		for e := range hooks {
			// You can handle additional events here if needed
			_ = e // Currently not used
		}
	}()

	myApp.Run()
}
