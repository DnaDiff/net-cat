package server

import "net"

type Client struct {
	username string
	remoteIP string
	conn     net.Conn
}

func addClient(clients *[]Client, username string, remoteIP string, conn net.Conn) {
	*clients = append(*clients, Client{username, remoteIP, conn})
}

func removeClient(clients *[]Client, remoteIP string) *[]Client {
	for i, client := range *clients {
		if client.remoteIP == remoteIP {
			*clients = append((*clients)[:i], (*clients)[i+1:]...)
		}
	}
	return clients
}
