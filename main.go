// TCP connection between server and multiple clients (relation of 1 to many).
// A name requirement to the client.
// Control connections quantity.
// Clients must be able to send messages to the chat.
// Do not broadcast EMPTY messages from a client.
// Messages sent, must be identified by the time that was sent and the user name of who sent the message, example : [2020-01-20 15:48:41][client.name]:[client.message]
// If a Client joins the chat, all the previous messages sent to the chat must be uploaded to the new Client.
// If a Client connects to the server, the rest of the Clients must be informed by the server that the Client joined the group.
// If a Client exits the chat, the rest of the Clients must be informed by the server that the Client left.
// All Clients must receive the messages sent by other Clients.
// If a Client leaves the chat, the rest of the Clients must not disconnect.
// If there is no port specified, then set as default the port 8989. Otherwise, program must respond with usage message: [USAGE]: ./TCPChat $port

package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8989"
	CONN_TYPE = "tcp"
)

// MessageQueue is a shared message queue among all clients.
type MessageQueue struct {
	head  *Node
	tail  *Node
	count int
}

// Node is a node in the linked list.
type Node struct {
	value string
	next  *Node
}

// AddMessage adds a new message to the queue.
func (q *MessageQueue) AddMessage(message string) {
	node := &Node{value: message}
	if q.head == nil {
		q.head = node
		q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.count++
}

// GetMessages gets all messages from the queue.
func (q *MessageQueue) GetMessages() []string {
	var messages []string
	current := q.head
	for current != nil {
		messages = append(messages, current.value)
		current = current.next
	}
	return messages
}

// RemoveMessage removes a message from the queue.
func (q *MessageQueue) RemoveMessage(message string) {
	if q.head == nil {
		return
	}
	if q.head.value == message {
		q.head = q.head.next
		q.count--
		return
	}
	current := q.head
	for current.next != nil {
		if current.next.value == message {
			current.next = current.next.next
			q.count--
			break
		}
		current = current.next
	}
	if q.head == nil {
		q.tail = nil
	}
}

func main() {
	// Create a new message queue
	messageQueue := MessageQueue{}

	// Listen for incoming connections
	ln, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes
	defer ln.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Accept incoming connections
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		// Handle connections in a new goroutine
		go handleConnection(conn, &messageQueue)
	}
}

// Handles incoming requests.
func handleConnection(conn net.Conn, messageQueue *MessageQueue) {
	// Send a welcome message to the client
	conn.Write([]byte("Welcome to the TCP Chat!\nUsername: "))
	fmt.Println("Incoming user...")
	// Make a buffer to hold incoming data
	buf := make([]byte, 1024)
	// Read the incoming data into the buffer
	rawUsername, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	username := strings.ReplaceAll(string(buf[:rawUsername]), "\n", "")
	// Send a response back to person contacting us
	conn.Write([]byte("\033[H\033[2JWe are glad to have you here, " + username + "!\nYou are now connected to the TCP Chat.\n"))
	fmt.Println("User '" + username + "' with IP address '" + conn.RemoteAddr().String() + "' connected to the TCP Chat.")

	// Start a goroutine that continuously sends messages from the queue to the client
	go func() {
		for {
			// Send the messages from the queue to the client
			for _, message := range messageQueue.GetMessages() {
				conn.Write([]byte(message + "\n"))
				messageQueue.RemoveMessage(message)
			}
			// Sleep for a short time before sending the messages again
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Listen for incoming messages
	for {
		message := receiveMessage(conn)
		if message == "exit" || message == "" {
			fmt.Println("User " + username + " disconnected from the TCP Chat.")
			conn.Close()
			break
		}
		// Add the message to the message queue
		messageQueue.AddMessage("[" + username + "]: " + message)
		fmt.Println("[" + username + "]: " + message)
	}
}

func receiveMessage(conn net.Conn) string {
	// Make a buffer to hold incoming data
	buf := make([]byte, 1024)
	// Read the incoming data into the buffer
	rawMessage, err := conn.Read(buf)
	if err != nil && err.Error() != "EOF" {
		fmt.Println("Error reading:", err.Error())
	}
	return strings.ReplaceAll(string(buf[:rawMessage]), "\n", "")
}
