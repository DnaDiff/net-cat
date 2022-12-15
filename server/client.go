package server

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	username string
	remoteIP string
	conn     net.Conn
}

type ClientList []Client

var mutex = &sync.Mutex{}

func clientHandler(clients *ClientList, messageLog *MessageLog, conn net.Conn) {

	if serverFull(clients, conn) {
		return
	}

	// Onboarding process
	pinguSender(conn, true)
	sendMessage(conn, MESSAGE_WELCOME)

	fmt.Println("Incoming user...")
	username := randomizeColor() + receiveMessage(conn) + "\033[0m"

	// Accept the user into the chat
	mutex.Lock()
	clients.AddClient(username, conn.RemoteAddr().String(), conn)
	mutex.Unlock()

	sendMessage(conn, fmt.Sprintf(MESSAGE_CONNECTED, username))

	// Print history for the user
	for _, message := range messageLog.Messages {
		sendMessage(conn, message+"\n")
	}

	fmt.Println("User '" + username + "' with IP address '" + conn.RemoteAddr().String() + "' connected to the TCP Chat.")
	clients.BroadcastMessage(messageLog, username, "\033[32mhas joined the chat.\033[0m")

	// Listen for incoming messages
	go func() {
		for {
			message := receiveMessage(conn)
			if message == "exit" {
				sendMessage(conn, MESSAGE_DISCONNECTED)
				pinguSender(conn, false)

				mutex.Lock()
				clients.RemoveClient(conn.RemoteAddr().String())
				mutex.Unlock()

				conn.Close()
				fmt.Println("User '" + username + "' disconnected from the TCP Chat.")
				clients.BroadcastMessage(messageLog, username, "\033[31mhas left the chat.\033[0m")
				break
			} else if message != "" {
				fmt.Printf(CHAT_FORMAT, getCurrentTime(), username, message+"\n")
				clients.BroadcastMessage(messageLog, username, message)
			}
		}
	}()
}

func (clients *ClientList) AddClient(username string, remoteIP string, conn net.Conn) {
	*clients = append(*clients, Client{username, remoteIP, conn})
}

func (clients *ClientList) RemoveClient(remoteIP string) {
	for i, client := range *clients {
		if client.remoteIP == remoteIP {
			*clients = append((*clients)[:i], (*clients)[i+1:]...)
		}
	}
}

func serverFull(clients *ClientList, conn net.Conn) bool {
	if len(*clients) >= MAX_CLIENTS {
		pinguSender(conn, false)
		sendMessage(conn, MESSAGE_FULL)
		conn.Close()
		return true
	}
	return false
}

func randomizeColor() string {
	rand.Seed(time.Now().UnixNano())
	return "\033[38;5;" + strconv.Itoa(rand.Intn(230)) + "m"
}
