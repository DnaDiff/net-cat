package main

import (
	"fmt"
	"net"
)

const (
	port = "8080"
	// maxMessageSize = 256
)

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
			fmt.Println(err)
			continue
		}
		// Handle the connection in a new goroutine
		go handleConnection(conn)
		// Send a message to the client
		sendMessage(conn, "Nicec!\n")
	}
}

func handleConnection(conn net.Conn) {
	// Do something with the connection
	fmt.Println("Received a new connection!")
	// Create a goroutine to read messages from the client
	// go func() {
	// 	for {
	// 		// Read data from the connection
	// 		data := make([]byte, maxMessageSize)
	// 		_, err := conn.Read(data)
	// 		if err != nil {
	// 			// Handle the error
	// 			fmt.Println(err)
	// 			return
	// 		}
	// 		// Print the received message
	// 		fmt.Println("Received message:", string(data))
	// 	}
	// }()
	// Don't forget to close the connection when you're done
	conn.Close()
}

func sendMessage(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		// Handle the error
		fmt.Println(err)
	}
}
