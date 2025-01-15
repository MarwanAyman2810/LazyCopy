// main.go
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

type clickableLabel struct {
	widget.Label
	lastClick time.Time
	onClick   func()
}

func newClickableLabel(text string, onClick func()) *clickableLabel {
	label := &clickableLabel{onClick: onClick}
	label.Text = text
	return label
}

func (c *clickableLabel) MouseDown(me *desktop.MouseEvent) {
	now := time.Now()
	if now.Sub(c.lastClick) < 500*time.Millisecond {
		if c.onClick != nil {
			c.onClick()
		}
	}
	c.lastClick = now
}

func (c *clickableLabel) MouseUp(*desktop.MouseEvent)    {}
func (c *clickableLabel) MouseIn(*desktop.MouseEvent)    {}
func (c *clickableLabel) MouseOut()                      {}
func (c *clickableLabel) MouseMoved(*desktop.MouseEvent) {}

func main() {
	myApp := app.New()

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

	mainWindow := myApp.NewWindow("LazyCopy")
	label := newClickableLabel("LazyCopy is running...\nDouble-click anywhere to show copy button", func() {
		copyWindow.CenterOnScreen()
		copyWindow.Show()
		copyWindow.RequestFocus()
	})

	mainWindow.SetContent(label)
	mainWindow.Resize(fyne.NewSize(200, 50))
	mainWindow.Show()

	myApp.Run()
}
