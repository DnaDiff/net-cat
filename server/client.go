package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	username    string
	remoteIP    string
	room        *Room
	lastMessage *string
	conn        net.Conn
}

var mutex = &sync.Mutex{}
var clientCount int

func InitializeClient(username string, remoteIP string, room *Room, conn net.Conn) *Client {
	client := &Client{username, remoteIP, room, new(string), conn}
	mutex.Lock()
	room.Clients = append(room.Clients, *client)
	mutex.Unlock()

	return client
}

func (client *Client) Disconnect() {
	client.room.RemoveClient(client.remoteIP)
	client.conn.Close()
	clientCount -= 1
}

func serverFull(conn net.Conn) bool {
	if clientCount >= MAX_CLIENTS {
		pinguSender(conn, false)
		SendMessage(conn, MESSAGE_ERROR_FULL)
		conn.Close()
		return true
	}
	return false
}

func inputListener(client *Client) {
	for {
		// Store the message
		*client.lastMessage = ReceiveMessage(client.conn)

		// Check if it is a command, run the command, break if it is an exit command
		// If it is not a command and the message is not empty, broadcast it
		if strings.HasPrefix(*client.lastMessage, "/") {
			if commandHandler(client, (*client.lastMessage)[1:]) {
				break
			}
		} else if *client.lastMessage != "" {
			client.room.AddMessage(client.username, *client.lastMessage)
		}
	}
}

func usernameCheck(conn net.Conn) (string, error) {
	username := ""
	for {
		username = ReceiveMessage(conn)
		// If the client was forcibly disconnected before receiving a name, don't name it /exit
		if strings.Contains(username, "/exit") {
			return "/exit", fmt.Errorf("Client terminated connection")
		}
		// If the username is three or more valid letters, pass
		if validNameBool(username) {
			username = randomizeColor() + username + "\033[0m"
			break
		} else {
			SendMessage(conn, MESSAGE_ERROR_USERNAME)
		}
	}
	return username, nil
}

func validNameBool(username string) bool {
	if len(username) >= 3 && strings.ContainsAny(strings.ToLower(username), "abcdefghijklmnopqrstuvwxyz") {
		return true
	}
	return false
}

func clientHandler(roomList map[string]*Room, room *Room, conn net.Conn) {

	if serverFull(conn) {
		return
	}

	// Onboarding process
	pinguSender(conn, true)
	SendMessage(conn, MESSAGE_CLIENT_WELCOME)

	clientCount += 1

	fmt.Fprintln(mw, "Incoming user...")
	username, err := usernameCheck(conn)
	if err != nil {
		clientCount -= 1

		conn.Close()
		fmt.Fprintln(mw, "User aborted")
		return
	}
	// Accept the user into the chat

	client := InitializeClient(username, conn.RemoteAddr().String(), room, conn)
	/* client := room.AddClient(username, conn.RemoteAddr().String(), "", conn) */

	// Clear the screen
	SendMessage(conn, MESSAGE_CLEAR)

	// Print history for the user
	for _, message := range room.Messages {
		SendMessage(conn, message+"\n")
	}

	SendMessage(conn, fmt.Sprintf(MESSAGE_CLIENT_CONNECTED, username, client.room.Name))

	fmt.Fprintln(mw, "User '"+username+"' with IP address '"+conn.RemoteAddr().String()+"' connected to the TCP Chat.")
	room.AddMessage(username, MESSAGE_ACTION_JOIN)

	// Listen for incoming messages
	inputListener(client)
}
