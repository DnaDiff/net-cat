package server

import (
	"fmt"
	"net"
)

type Client struct {
	username string
	remoteIP string
	conn     net.Conn
}

type ClientList []Client

func clientHandler(clients *ClientList, messageLog *MessageLog, conn net.Conn) {
	// Onboarding process
	PinguSender(conn, true)
	sendMessage(conn, welcomeMessage)

	fmt.Println("Incoming user...")
	username := receiveMessage(conn)

	// Accept the user into the chat
	clients.addClient(username, conn.RemoteAddr().String(), conn)

	sendMessage(conn, fmt.Sprintf(connectedMessage, username))

	fmt.Println("User '" + username + "' with IP address '" + conn.RemoteAddr().String() + "' connected to the TCP Chat.")

	// Listen for incoming messages
	go func() {
		// Print history for the user
		for _, message := range messageLog.Messages {
			sendMessage(conn, message+"\n")
		}

		for {
			message := receiveMessage(conn)
			if message == "exit" {
				fmt.Println("User " + username + " disconnected from the TCP Chat.")
				sendMessage(conn, disconnectedMessage)
				PinguSender(conn, false)

				clients.removeClient(conn.RemoteAddr().String())
				conn.Close()
				break
			} else {
				broadcastMessage(clients, messageLog, username, message)
			}
			fmt.Printf(CHAT_FORMAT, getCurrentTime(), username, message+"\n")
		}
	}()
}

func (clients *ClientList) addClient(username string, remoteIP string, conn net.Conn) {
	*clients = append(*clients, Client{username, remoteIP, conn})
}

func (clients *ClientList) removeClient(remoteIP string) {
	for i, client := range *clients {
		if client.remoteIP == remoteIP {
			*clients = append((*clients)[:i], (*clients)[i+1:]...)
		}
	}
}
