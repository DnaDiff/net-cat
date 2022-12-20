package server

import (
	"net"
	"time"
)

// Message templates for the client

const (
	MESSAGE_CLEAR = "\033[H\033[2J"

	MESSAGE_CLIENT_WELCOME      = "Welcome to the TCP Chat!\nUsername: "
	MESSAGE_CLIENT_CONNECTED    = "We are glad to have you here, %s! You are now connected to room %s.\nList other rooms with \033[33m/rooms\033[0m. To disconnect, type \033[33m/exit\033[0m.\n"
	MESSAGE_CLIENT_DISCONNECTED = "You have been disconnected from the TCP Chat.\n\033[3mPingu is sad to see you go :(\nPingu will miss you!\nPingu will cry!\nPingu will die!\n\033[0mPress [ENTER] to leave.\n"
	MESSAGE_CLIENT_SWITCH       = "You are now chatting in room %s.\nReturn by typing \033[33m/room %s.\n"

	MESSAGE_ACTION_JOIN  = "\033[32mhas joined the server\033[0m"
	MESSAGE_ACTION_LEAVE = "\033[31mhas left the server\033[0m"
	MESSAGE_ACTION_NAME  = "\033[33mis now known as \033[0m%s"

	MESSAGE_ROOM_JOIN  = "\033[33mjoined the room\033[0m"
	MESSAGE_ROOM_LEAVE = "\033[33mwent to room \033[36m%s\033[33m\033[0m"

	MESSAGE_COMMAND_HELP  = "\033[33m/%s \033[1;30m- \033[3m%s\033[0m\n"
	MESSAGE_COMMAND_ROOMS = "%s - \033[33m%d\033[0m user(s)\n"

	MESSAGE_ERROR_FULL       = "Pingu is sad to tell you that the chat is full. Please come back to play with Pingu at a later time.\nPress [ENTER] to leave.\n"
	MESSAGE_ERROR_USERNAME   = "\033[3mYour username has to contain at least three valid letters.\033[0m\nUsername: "
	MESSAGE_ERROR_HELP       = "Type \033[33m\033[3m/help\033[0m to see available commands.\033[0m\n"
	MESSAGE_ERROR_ROOM       = "\033[3mYou have already joined this room.\033[0m\n"
	MESSAGE_ERROR_USAGE_ROOM = "Usage: \033[3m/room <room_name>\033[0m\n"
)

var pinguAlive = []string{
	"         _nnnn_",
	"        dGGGGMMb",
	"       @p~qp~~qMb",
	"       M|@||@) M|",
	"       @,----.JM|",
	"      JS^\\__/  qKL",
	"     dZP        qKRb",
	"    dZP          qKKb",
	"   fZP    helo    SMMb",
	"   HZM            MMMM",
	"   FqM            MMMM",
	" __| \".        |\\dS\"qML",
	" |    `.       | `' \\Zq",
	"_)      \\.___.,|     .'",
	"\\____   )MMMMMP|   .'",
	"     `-'       `--'",
}

var pinguDead = []string{
	"         _nnnn_",
	"        dGGGGMMb",
	"       @p~qp~~qMb",
	"       M|X||X) M|",
	"       @,----.JM|",
	"      JS^\\__/  qKL",
	"     dZP        qKRb",
	"    dZP          qKKb",
	"   fZP    goby    SMMb",
	"   HZM            MMMM",
	"   FqM           MMMM",
	" __| \".  𒍼      |\\dS\"qML",
	" |    `.        | `' \\Zq",
	"_)      \\.___.,|     .'",
	"\\____   )MMMMMP|   .'",
	"     `-'       `--'",
}

func pinguSender(conn net.Conn, isAlive bool) {
	if isAlive {
		for _, e := range pinguAlive {
			SendMessage(conn, e+"\n")
		}
	} else {
		for _, e := range pinguDead {
			SendMessage(conn, e+"\n")
		}
	}
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
