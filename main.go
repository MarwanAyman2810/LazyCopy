package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

// ClickCatcher implements desktop.Mouseable
type ClickCatcher struct {
	widget.BaseWidget
	lastClick  time.Time
	copyWindow fyne.Window
}

func NewClickCatcher(copyWin fyne.Window) *ClickCatcher {
	c := &ClickCatcher{copyWindow: copyWin}
	c.ExtendBaseWidget(c)
	return c
}

func (c *ClickCatcher) MouseDown(*desktop.MouseEvent) {
	fmt.Println("Mouse click detected!")
	now := time.Now()
	if now.Sub(c.lastClick) < 500*time.Millisecond {
		fmt.Println("Double click detected! Showing window...")
		c.copyWindow.CenterOnScreen()
		c.copyWindow.Show()
		c.copyWindow.RequestFocus()
	}
	c.lastClick = now
}

func (c *ClickCatcher) MouseUp(*desktop.MouseEvent) {}
func (c *ClickCatcher) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewPadded())
}

func main() {
	fmt.Println("Starting application...")
	myApp := app.New()

	// Create only the copy button window
	fmt.Println("Creating copy window...")
	copyWindow := myApp.NewWindow("Copy")
	copyWindow.Resize(fyne.NewSize(100, 40))
	copyWindow.SetFixedSize(true)
	copyWindow.Hide()

	copyBtn := widget.NewButton("Copy", func() {
		fmt.Println("Copy button clicked!")
		if err := clipboard.WriteAll("Hello from the app!"); err != nil {
			fmt.Printf("Error writing to clipboard: %v\n", err)
		}
		copyWindow.Hide()
	})
	copyWindow.SetContent(container.NewPadded(copyBtn))

	// Create a visible main window
	mainWindow := myApp.NewWindow("Click Catcher")
	mainWindow.Resize(fyne.NewSize(200, 100))

	clickCatcher := NewClickCatcher(copyWindow)
	mainWindow.SetContent(clickCatcher)

	mainWindow.Show()
	fmt.Println("Starting main application loop...")
	myApp.Run()
}
