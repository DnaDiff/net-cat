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

type ClientList struct {
	clients []Client
	count   int
}

var mutex = &sync.Mutex{}

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
	username := randomizeColor() + receiveMessage(conn) + "\033[0m"

	// Accept the user into the chat

	clientList.AddClient(username, conn.RemoteAddr().String(), conn)

	sendMessage(conn, fmt.Sprintf(MESSAGE_CONNECTED, username))

	// Print history for the user
	for _, message := range messageLog.Messages {
		sendMessage(conn, message+"\n")
	}

	fmt.Println("User '" + username + "' with IP address '" + conn.RemoteAddr().String() + "' connected to the TCP Chat.")
	clientList.BroadcastMessage(messageLog, username, "\033[32mhas joined the chat.\033[0m")

	// Listen for incoming messages
	go func() {
		for {
			message := receiveMessage(conn)
			if message == "exit" {
				sendMessage(conn, MESSAGE_DISCONNECTED)
				pinguSender(conn, false)

				clientList.RemoveClient(conn.RemoteAddr().String())

				conn.Close()
				fmt.Println("User '" + username + "' disconnected from the TCP Chat.")
				clientList.BroadcastMessage(messageLog, username, "\033[31mhas left the chat.\033[0m")
				break
			} else if message != "" {
				fmt.Printf(CHAT_FORMAT, getCurrentTime(), username, message+"\n")
				clientList.BroadcastMessage(messageLog, username, message)
			}
		}
	}()
}

func (clientList *ClientList) AddClient(username string, remoteIP string, conn net.Conn) {
	mutex.Lock()
	fmt.Println(*clientList)
	clientList.clients = append(clientList.clients, Client{username, remoteIP, conn})
	// clientList.count++ // We normally increment the count before AddClient is run, to catch partial connections
	mutex.Unlock()

	fmt.Println(*clientList)
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

func randomizeColor() string {
	rand.Seed(time.Now().UnixNano())
	return "\033[38;5;" + strconv.Itoa(rand.Intn(230)) + "m"
}
