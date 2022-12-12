package server

import (
	"net"
	"time"
)

// Message templates for the client

const (
	welcomeMessage      = "Welcome to the TCP Chat!\nUsername: "
	connectedMessage    = "\033[H\033[2JWe are glad to have you here, %s!\nYou are now connected to the TCP Chat.\n"
	disconnectedMessage = "You have been disconnected from the TCP Chat.\nPingu is sad to see you go :(\nPingu will miss you!\nPingu will cry!\nPingu will die!\nPress [ENTER] to leave.\n"
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
	" __| \".  𒍼   |\\dS\"qML",
	" |    `.        | `' \\Zq",
	"_)      \\.___.,|     .'",
	"\\____   )MMMMMP|   .'",
	"     `-'       `--'",
}

func pinguSender(conn net.Conn, isAlive bool) {
	if isAlive {
		for _, e := range pinguAlive {
			sendMessage(conn, e+"\n")
		}
	} else {
		for _, e := range pinguDead {
			sendMessage(conn, e+"\n")
		}
	}
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
