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
	lastMessage string
	conn        net.Conn
}

var mutex = &sync.Mutex{}
var clientCount int

func InitializeClient(username string, remoteIP string, room *Room, lastMessage string, conn net.Conn) *Client {
	client := &Client{username, remoteIP, room, lastMessage, conn}
	mutex.Lock()
	room.Clients = append(room.Clients, *client)
	mutex.Unlock()

	return client
}

func (room *Room) AddClient(client *Client) {
	mutex.Lock()
	room.Clients = append(room.Clients, *client)
	mutex.Unlock()
}

func (room *Room) RemoveClient(remoteIP string) {
	mutex.Lock()
	for i, client := range room.Clients {
		if client.remoteIP == remoteIP {
			room.Clients = append((room.Clients)[:i], (room.Clients)[i+1:]...)
		}
	}
	mutex.Unlock()
}

func CreateRoom(roomList map[string]*Room, roomName string) {
	roomExists := false
	for room := range roomList {
		if room == roomName {
			roomExists = true
			break
		}
	}
	if !roomExists {
		roomList[roomName] = &Room{
			Name:     roomName,
			Clients:  []Client{},
			Messages: []string{},
		}
	}
}

func (client *Client) SwitchRoom(roomList map[string]*Room, roomName string) {
	CreateRoom(roomList, roomName)
	client.room.RemoveClient(client.remoteIP)
	roomList[roomName].AddClient(client)
}

func (client *Client) Disconnect() {
	client.room.RemoveClient(client.remoteIP)
	client.conn.Close()
	clientCount -= 1
}

func serverFull(conn net.Conn) bool {
	if clientCount >= MAX_CLIENTS {
		pinguSender(conn, false)
		sendMessage(conn, MESSAGE_FULL)
		conn.Close()
		return true
	}
	return false
}

func inputListener(client *Client) {
	for {
		// Store the message
		client.lastMessage = receiveMessage(client.conn)

		// Check if it is a command, run the command, break if it is an exit command
		// If it is not a command and the message is not empty, broadcast it
		if strings.HasPrefix(client.lastMessage, "/") {
			if commandHandler(client, client.lastMessage[1:]) {
				break
			}
		} else if client.lastMessage != "" {
			client.room.AddMessage(client.username, client.lastMessage)
		}
	}
}

func usernameCheck(conn net.Conn) (string, error) {
	username := ""
	for {
		username = receiveMessage(conn)
		// If the client was forcibly disconnected before receiving a name, don't name it /exit
		if strings.Contains(username, "/exit") {
			return "/exit", fmt.Errorf("Client terminated connection")
		}
		// If the username is three or more valid letters, pass
		if validNameBool(username) {
			username = randomizeColor() + username + "\033[0m"
			break
		} else {
			sendMessage(conn, MESSAGE_USERNAME_ERROR)
		}
	}
	return username, nil
}

func validNameBool(username string) bool {
	if len(username) >= 3 && strings.ContainsAny(username, "abcdefghijklmnopqrstuvwxyz") {
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
	sendMessage(conn, MESSAGE_WELCOME)

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

	client := InitializeClient(username, conn.RemoteAddr().String(), room, "", conn)
	/* client := room.AddClient(username, conn.RemoteAddr().String(), "", conn) */

	sendMessage(conn, fmt.Sprintf(MESSAGE_CONNECTED, username))

	// Print history for the user
	for _, message := range room.Messages {
		sendMessage(conn, message+"\n")
	}

	fmt.Fprintln(mw, "User '"+username+"' with IP address '"+conn.RemoteAddr().String()+"' connected to the TCP Chat.")
	room.AddMessage(username, "\033[32mhas joined the chat.\033[0m")

	// Listen for incoming messages
	inputListener(client)
}
