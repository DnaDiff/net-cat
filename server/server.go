package server

import (
	"fmt"
	"net"
	"os"
)

func StartServer(CONN_PORT string) {
	var clients ClientList
	var messageLog MessageLog

	ln, err := net.Listen("tcp", ":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes
	defer ln.Close()

	fmt.Println("Listening on " + ":" + CONN_PORT)
	for {
		// Accept incoming connections
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		// Handle connections in a new goroutine
		go clientHandler(&clients, &messageLog, conn)
	}
}
