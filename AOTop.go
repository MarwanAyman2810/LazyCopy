package main

import (
	"encoding/binary"
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func InitializeAOTopWindow() {
	// Connect to the X server
	conn, err := xgb.NewConn()
	if err != nil {
		log.Fatalf("Failed to connect to X server: %v", err)
	}
	defer conn.Close()

	// Get the default screen
	setup := xproto.Setup(conn)
	screen := setup.DefaultScreen(conn)

	// Create a simple window
	window, err := xproto.NewWindowId(conn)
	if err != nil {
		log.Fatalf("Failed to create window ID: %v", err)
	}

	xproto.CreateWindow(
		conn,
		screen.RootDepth,
		window,
		screen.Root,
		0, 0, 300, 200,
		0,
		xproto.WindowClassInputOutput,
		screen.RootVisual,
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{screen.WhitePixel, xproto.EventMaskExposure | xproto.EventMaskKeyPress},
	)

	// Set the window title
	title := "Always On Top Window"
	xproto.ChangeProperty(
		conn,
		xproto.PropModeReplace,
		window,
		xproto.AtomWmName,
		xproto.AtomString,
		8,
		uint32(len(title)),
		[]byte(title),
	)

	// Map (show) the window
	xproto.MapWindow(conn, window)

	// Set the window to be always on top
	setAlwaysOnTop(conn, window)

	// Event loop
	for {
		ev, err := conn.WaitForEvent()
		if err != nil {
			log.Fatalf("Error waiting for event: %v", err)
		}
		switch ev.(type) {
		case xproto.ExposeEvent:
			// Handle expose events if necessary
		case xproto.KeyPressEvent:
			// Exit on key press
			return
		}
	}
}

func setAlwaysOnTop(conn *xgb.Conn, window xproto.Window) {
	// Get the atom for _NET_WM_STATE
	atomName := "_NET_WM_STATE"
	wmStateCookie := xproto.InternAtom(conn, true, uint16(len(atomName)), atomName)
	wmStateReply, err := wmStateCookie.Reply()
	if err != nil {
		log.Fatalf("Failed to get atom for _NET_WM_STATE: %v", err)
	}
	wmState := wmStateReply.Atom

	// Get the atom for _NET_WM_STATE_ABOVE
	atomName = "_NET_WM_STATE_ABOVE"
	wmStateAboveCookie := xproto.InternAtom(conn, true, uint16(len(atomName)), atomName)
	wmStateAboveReply, err := wmStateAboveCookie.Reply()
	if err != nil {
		log.Fatalf("Failed to get atom for _NET_WM_STATE_ABOVE: %v", err)
	}
	wmStateAbove := wmStateAboveReply.Atom

	// Prepare the data to set the window property
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(wmStateAbove))

	// Change the window property to set it as always on top
	xproto.ChangeProperty(
		conn,
		xproto.PropModeReplace,
		window,
		wmState,
		xproto.AtomAtom,
		32,
		1,
		data,
	)
}
