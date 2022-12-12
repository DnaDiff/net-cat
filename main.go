// TCP connection between server and multiple clients (relation of 1 to many).
// A name requirement to the client.
// Control connections quantity.
// Clients must be able to send messages to the chat.
// Do not broadcast EMPTY messages from a client.
// Messages sent, must be identified by the time that was sent and the user name of who sent the message, example : [2020-01-20 15:48:41][client.name]:[client.message]
// If a Client joins the chat, all the previous messages sent to the chat must be uploaded to the new Client.
// If a Client connects to the server, the rest of the Clients must be informed by the server that the Client joined the group.
// If a Client exits the chat, the rest of the Clients must be informed by the server that the Client left.
// All Clients must receive the messages sent by other Clients.
// If a Client leaves the chat, the rest of the Clients must not disconnect.
// If there is no port specified, then set as default the port 8989. Otherwise, program must respond with usage message: [USAGE]: ./TCPChat $port

package main

import (
	"fmt"
	"os"

	"github.com/DnaDiff/net-cat/server"
)

var CONN_PORT = "8989"

// TODO: Sort functions into their own categories
// TODO: Fix syncing:
// - When a client connects, send all previous messages to the client
// - When a client disconnects, send a message to all other clients that the client disconnected
// - When a client sends a message, send the message to all other clients:
//   - The message should be in the format: [time][username]:[message]
// - Syncing should be done in a separate goroutine

func main() {
	if len(os.Args) == 2 {
		CONN_PORT = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
	server.StartServer(CONN_PORT)
}
