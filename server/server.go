package server

import (
	"fmt"
	"net"
	"os"
)

const (
	CONN_TYPE = "tcp"

	MAX_CLIENTS = 10
)

var isRunning bool = false
var ln net.Listener

func StartServer(port string, logFlag bool) error {
	var roomList = map[string]*Room{
		"general": {
			Name:     "general",
			Clients:  []Client{},
			Messages: []string{},
		},
	}

	logger, err := logCheck(logFlag)
	if err != nil {
		logger.Close()
		return err
	}
	defer logger.Close()

	ln, err = net.Listen(CONN_TYPE, ":"+port)
	if err != nil {
		return err
	}
	// Close the listener when the application closes
	defer ln.Close()

	fmt.Fprintln(mw, "Listening on "+":"+port)
	isRunning = true
	for isRunning {
		// Pause the loop and wait for incoming connections to accept
		conn, err := ln.Accept()
		if err != nil {
			// Prevent shutdown-related errors
			if !isRunning {
				break
			}
			fmt.Fprintln(mw, "Error accepting: ", err.Error())
			continue
		}
		// Handle connections in a new goroutine
		go clientHandler(roomList, roomList["general"], conn)
	}
	fmt.Fprintln(mw, "Server stopped")
	return nil
}

func StopServer() {
	isRunning = false
	ln.Close()
}

func logCheck(flag bool) (os.File, error) {
	var logger os.File
	var err error
	if flag {
		logger, err = enableLogging(true)
		if err != nil {
			return logger, err
		}
	} else {
		logger, err = enableLogging(false)
		if err != nil {
			return logger, err
		}
	}
	return logger, nil
}
