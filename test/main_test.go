package main

import (
	"fmt"
	"io"
	"net"
	"os/exec"
	"testing"
	"time"

	"github.com/DnaDiff/net-cat/server"
)

const CONN_PORT = "8989"

func TestConnections(t *testing.T) {
	go TestServer(t)
	for i := 0; i < 10; i++ {
		go TestClient(t)
	}
	time.Sleep(10 * time.Second)
}

// Test server
func TestServer(t *testing.T) {
	// Start the server in a goroutine
	go func() {
		err := server.StartServer(CONN_PORT, false)
		if err != nil {
			t.Errorf("Error starting server: %v", err)
		}
	}()

	// Wait for next action
	time.Sleep(10 * time.Second)

	// Stop the server
	server.StopServer()

	// Check if the server is still running on the specified port
	conn, err := net.Dial("tcp", "localhost:"+CONN_PORT)
	if err == nil {
		t.Errorf("Expected the server to be stopped, but it is still running")
		conn.Close()
	}

	time.Sleep(5 * time.Second)
}

// Test case for user connecting, sending a message, and disconnecting
func TestClient(t *testing.T) {
	// Check if the server is running on the specified port
	go func() {
		_, err := net.Dial("tcp", "localhost:"+CONN_PORT)
		if err != nil {
			TestServer(t)
		}
	}()

	// Create a new exec.Cmd struct to run the "nc" command
	cmd := exec.Command("nc", "localhost", "8989")

	// Use the Stdin pipe of the cmd struct to write to the command's standard input
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Errorf("Error creating StdinPipe for Cmd: %v", err)
	}

	// Start the command in a new goroutine so that it doesn't block the current thread
	go func() {
		sendMessage(t, stdin, "Bot")
		time.Sleep(100 * time.Millisecond)
		sendMessage(t, stdin, "Hello, world!")
		time.Sleep(100 * time.Millisecond)
		sendMessage(t, stdin, "exit")
		time.Sleep(100 * time.Millisecond)

		// Close the Stdin pipe to signal to the "nc" command that we're done writing to it
		if err := stdin.Close(); err != nil {
			t.Errorf("Error closing StdinPipe for Cmd: %v", err)
		}
	}()

	// Run the command and wait for it to finish
	err = cmd.Run()
	if err != nil {
		if err.Error() == "exit status 1" {
			// Connection was closed before sending messages
			t.Errorf("Connection closed before sending messages: %v", err)
		} else {
			// Other error occurred
			t.Errorf("Error running command: %v", err)
		}
	}
	time.Sleep(8 * time.Second)
}

func sendMessage(t *testing.T, stdin io.WriteCloser, message string) {
	// Write the desired input to the command's standard input
	_, err := fmt.Fprintln(stdin, message)
	if err != nil {
		t.Errorf("Error writing message: %v, %v", err, message)
	}

	// Prevent the client from sending messages too quickly
	time.Sleep(100 * time.Millisecond)
}
