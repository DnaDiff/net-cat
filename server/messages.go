package server

import (
	"fmt"
	"net"
	"strings"
)

const CHAT_FORMAT = "[%s][%s]: %s"

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
	conn.Write([]byte(message))
}

func broadcastMessage(clients *ClientList, messageLog *MessageLog, messageUsername string, message string) {
	messageLog.AddMessage(fmt.Sprintf(CHAT_FORMAT, getCurrentTime(), messageUsername, message))
	for _, client := range *clients {
		if client.username == messageUsername {
			sendMessage(client.conn, "\r\033[1A\033[2K")
		}
		sendMessage(client.conn, fmt.Sprintf(CHAT_FORMAT, getCurrentTime(), messageUsername, message+"\n"))
	}
}
