package server

import (
	"fmt"
	"net"
)

func handleClient(clients *[]Client, messageLog *MessageLog, conn net.Conn) {
	// Onboarding process
	conn.Write([]byte(welcomeMessage))

	fmt.Println("Incoming user...")
	username := receiveMessage(conn)

	// Accept the user into the chat
	addClient(clients, username, conn.RemoteAddr().String(), conn)

	conn.Write([]byte(connectedMessage))

	fmt.Println("User '" + username + "' with IP address '" + conn.RemoteAddr().String() + "' connected to the TCP Chat.")

	// Listen for incoming messages
	go func() {
		// Print history for the user
		for _, message := range messageLog.Messages {
			sendMessage(conn, message)
		}

		for {
			message := receiveMessage(conn)
			if message == "exit" {
				fmt.Println("User " + username + " disconnected from the TCP Chat.")
				clients = removeClient(clients, conn.RemoteAddr().String())
				conn.Close()
				break
			} else {
				broadcastMessage(clients, messageLog, "["+username+"]:"+message)
			}
			fmt.Println("[" + username + "]: " + message)
		}
	}()
}
