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

	if len(*clients) >= 10 {
		pinguSender(conn, false)
		sendMessage(conn, "Pingu is sad to tell you that the chat is full. Please come back to play with Pingu at a later time.")
		conn.Close()
		return
	}

	// Onboarding process
	pinguSender(conn, true)
	sendMessage(conn, welcomeMessage)

	fmt.Println("Incoming user...")
	username := randomizeColor() + receiveMessage(conn) + "\033[0m"

	// Accept the user into the chat
	mutex.Lock()
	clients.AddClient(username, conn.RemoteAddr().String(), conn)
	mutex.Unlock()

	sendMessage(conn, fmt.Sprintf(connectedMessage, username))

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
				sendMessage(conn, disconnectedMessage)
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

func randomizeColor() string {
	rand.Seed(time.Now().UnixNano())
	return "\033[38;5;" + strconv.Itoa(rand.Intn(230)) + "m"
}
