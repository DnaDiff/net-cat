package server

import (
	"fmt"
	"net"
	"strings"
)

// MessageQueue is a shared message queue among all clients.
type MessageLog struct {
	Messages []string
	count    int
}

// AddMessage adds a new message to the queue.
func (q *MessageLog) AddMessage(message string) {
	q.Messages = append(q.Messages, message)
	q.count++
}

func receiveMessage(conn net.Conn) string {
	// Make a buffer to hold incoming data
	buf := make([]byte, 1024)
	// Read the incoming data into the buffer
	rawMessage, err := conn.Read(buf)
	if err != nil && rawMessage != 0 {
		fmt.Println("Error reading:", err.Error())
	} else if rawMessage == 0 {
		return "exit"
	}
	return strings.ReplaceAll(string(buf[:rawMessage]), "\n", "")
}

func sendMessage(conn net.Conn, message string) {
	conn.Write([]byte(message + "\n"))
}

func broadcastMessage(clients *[]Client, messageLog *MessageLog, message string) {
	messageLog.AddMessage(message)
	for _, client := range *clients {
		sendMessage(client.conn, message)
	}
}
