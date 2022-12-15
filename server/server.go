package server

import (
	"fmt"
	"net"
)

const (
	CONN_TYPE = "tcp"

	MAX_CLIENTS = 10
)

var isRunning bool = false
var ln net.Listener

func StartServer(CONN_PORT string) error {
	var clients ClientList
	var messageLog MessageLog

	var err error // Prevent shadowing of ln below
	ln, err = net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		return err
	}
	// Close the listener when the application closes
	defer ln.Close()

	fmt.Println("Listening on " + ":" + CONN_PORT)
	isRunning = true
	for isRunning {
		// Pause the loop and wait for incoming connections to accept
		conn, err := ln.Accept()
		if err != nil {
			// Prevent shutdown-related errors
			if !isRunning {
				break
			}
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		// Handle connections in a new goroutine
		go clientHandler(&clients, &messageLog, conn)
	}
	fmt.Println("Server stopped")
	return nil
}

func StopServer() {
	isRunning = false
	ln.Close()
}
