package server

import (
	"fmt"
	"io"
	"net"
	"strings"
)

const CHAT_FORMAT = "[%s][%s]: %s" // [date + time][username]: [message]

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
	} else if err == io.EOF {
		return "/exit"
	}
	sendMessage(conn, "\r\033[1A\033[2K")
	return strings.ReplaceAll(string(buf[:rawMessage]), "\n", "")
}

func sendMessage(conn net.Conn, message string) {
	conn.Write([]byte(message))
}

func (clientList *ClientList) BroadcastMessage(messageLog *MessageLog, messageUsername string, message string) {
	fmt.Println(fmt.Sprintf(CHAT_FORMAT, getCurrentTime(), messageUsername, message))
	mutex.Lock()
	messageLog.AddMessage(fmt.Sprintf(CHAT_FORMAT, getCurrentTime(), messageUsername, message))
	mutex.Unlock()
	for _, client := range clientList.clients {
		sendMessage(client.conn, fmt.Sprintf(CHAT_FORMAT, getCurrentTime(), messageUsername, message+"\n"))
		// sendMessage(client.conn, client.username+"> ")
	}
}
