package main

import (
	"fmt"
	"net"
)

var port = "8080"

func main() {
	// Create a listener on port 8080
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// Handle the error
		panic(err)
	}
	fmt.Println("Listening on port " + port)
	defer ln.Close()

	// Accept incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			// Handle the error
			continue
		}
		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Do something with the connection
	fmt.Println("Received a new connection!")
	// Don't forget to close the connection when you're done
	conn.Close()
}
