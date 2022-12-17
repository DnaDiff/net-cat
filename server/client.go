package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	username string
	remoteIP string
	conn     net.Conn
}

type ClientList struct {
	clients []Client
	count   int
}

var mutex = &sync.Mutex{}

func (clientList *ClientList) AddClient(username string, remoteIP string, conn net.Conn) *Client {
	client := &Client{username, remoteIP, conn}
	mutex.Lock()
	fmt.Println(*clientList)
	clientList.clients = append(clientList.clients, *client)
	// clientList.count++ // We normally increment the count before AddClient is run, to catch partial connections
	mutex.Unlock()

	fmt.Println(*clientList)
	return client
}

func (clientList *ClientList) RemoveClient(remoteIP string) {
	mutex.Lock()
	for i, client := range clientList.clients {
		if client.remoteIP == remoteIP {
			clientList.clients = append((clientList.clients)[:i], (clientList.clients)[i+1:]...)
		}
	}
	clientList.count--
	mutex.Unlock()
}

func serverFull(clientList *ClientList, conn net.Conn) bool {
	if clientList.count >= MAX_CLIENTS {
		pinguSender(conn, false)
		sendMessage(conn, MESSAGE_FULL)
		conn.Close()
		return true
	}
	return false
}

func inputListener(clientList *ClientList, client *Client, messageLog *MessageLog) {
	for {
		// Store the message
		message := receiveMessage(client.conn)

		// Check if it is a command, run the command, break if it is an exit command
		// If it is not a command and the message is not empty, broadcast it
		if strings.HasPrefix(message, "/") {
			if commandHandler(clientList, client, messageLog, message[1:]) {
				break
			}
		} else if message != "" {
			clientList.BroadcastMessage(messageLog, client.username, message)
		}
	}
}

func usernameCheck(conn net.Conn) string {
	username := ""
	for {
		username = receiveMessage(conn)
		if len(username) >= 3 && strings.ContainsAny(username, "abcdefghijklmnopqrstuvwxyz") {
			username = randomizeColor() + username + "\033[0m"
			break
		} else {
			sendMessage(conn, MESSAGE_USERNAME_ERROR)
		}
	}
	return username
}

func clientHandler(clientList *ClientList, messageLog *MessageLog, conn net.Conn) {

	if serverFull(clientList, conn) {
		return
	}

	// Onboarding process
	pinguSender(conn, true)
	sendMessage(conn, MESSAGE_WELCOME)

	mutex.Lock()
	clientList.count++
	mutex.Unlock()

	fmt.Println("Incoming user...")
	username := usernameCheck(conn)
	// Accept the user into the chat

	client := clientList.AddClient(username, conn.RemoteAddr().String(), conn)

	sendMessage(conn, fmt.Sprintf(MESSAGE_CONNECTED, username))

	// Print history for the user
	for _, message := range messageLog.Messages {
		sendMessage(conn, message+"\n")
	}

	fmt.Println("User '" + username + "' with IP address '" + conn.RemoteAddr().String() + "' connected to the TCP Chat.")
	clientList.BroadcastMessage(messageLog, username, "\033[32mhas joined the chat.\033[0m")

	// Listen for incoming messages
	go inputListener(clientList, client, messageLog)
}
