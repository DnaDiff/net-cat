package server

import (
	"fmt"
	"io"
	"net"
	"strings"
)

const CHAT_FORMAT = "[%s][%s]: %s" // [date + time][username]: [message]

type Room struct {
	Name     string
	Clients  []Client
	Messages []string
}

// AddMessage prints the message to the terminal, logs the message and sends it to all clients.
func (room *Room) AddMessage(messageUsername, message string) {
	fmt.Fprintln(mw, fmt.Sprintf("[%v]"+CHAT_FORMAT, room.Name, getCurrentTime(), messageUsername, message))
	mutex.Lock()
	room.Messages = append(room.Messages, fmt.Sprintf(CHAT_FORMAT, getCurrentTime(), messageUsername, message))
	mutex.Unlock()
	fmt.Fprintln(mw, room.Clients)
	for _, client := range room.Clients {
		sendMessage(client.conn, fmt.Sprintf(CHAT_FORMAT, getCurrentTime(), messageUsername, message+"\n"))
	}
}

// It reads the incoming data from the connection, and returns it as a string
func receiveMessage(conn net.Conn) string {
	// Make a buffer to hold incoming data
	buf := make([]byte, 1024)
	// Read the incoming data into the buffer
	rawMessage, err := conn.Read(buf)
	if err != nil && rawMessage != 0 {
		fmt.Fprintln(mw, "Error reading:", err.Error())
	} else if err == io.EOF {
		return "/exit"
	}
	// Remove client-side input line appended after pressing enter
	sendMessage(conn, "\r\033[1A\033[2K")
	return strings.ReplaceAll(string(buf[:rawMessage]), "\n", "")
}

// SendMessage writes the message to the connection.
func sendMessage(conn net.Conn, message string) {
	conn.Write([]byte(message))
}
