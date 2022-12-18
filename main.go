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
	"log"
	"os"
	"strings"

	"github.com/DnaDiff/net-cat/server"
)

var port = "8989"
var logFlag = false

func main() {
	getArgs()
	if err := server.StartServer(port, logFlag); err != nil {
		log.Fatal(err)
	}
}

func getArgs() {
	if strings.Contains(strings.Join(os.Args, " "), "--log") {
		logFlag = true
		fmt.Println("Logging enabled")
	}

	if len(os.Args) == 2 && !logFlag {
		port = os.Args[1]
	} else if len(os.Args) >= 3 && !logFlag {
		log.Fatal("[USAGE]: ./TCPChat $port [--log]")
	}
}
