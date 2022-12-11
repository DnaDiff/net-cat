package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	port           = "8080"
	maxMessageSize = 256
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
		// Send the pingu image to the client
		PinguSender(conn)
		// sendMessage(conn, "Pingu.txt\n")
		sendMessage(conn, "NOOOT NOOOT!\n")
	}
}

// Send the pingu image to the client
func PinguSender(conn net.Conn) {
	// get the contents of the pingu.txt file just like in go-reloaded
	file, err := os.Open("../pingu.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sendMessage(conn, scanner.Text()+"\n")
	}
	// check for errors
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func handleConnection(conn net.Conn) {
	// Do something with the connection
	fmt.Println("Received a new connection!")
	// Create a goroutine to read messages from the client
	go func() {
		for {
			readMessage(conn)
		}
	}()
	// conn.Close()
}

func readMessage(conn net.Conn) {
	// Read data from the connection
	data := make([]byte, maxMessageSize)
	_, err := conn.Read(data)
	if err != nil {
		// Handle the error
		fmt.Println(err)
		return
	}
	fmt.Println("Received message:", string(data))
}
func sendMessage(conn net.Conn, msg string) {
	// Write data to the connection
	_, err := conn.Write([]byte(msg))
	if err != nil {
		// Handle the error
		fmt.Println(err)
	}
}
