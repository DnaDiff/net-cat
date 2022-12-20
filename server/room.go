package server

import (
	"fmt"
	"strings"
)

func (room *Room) AddClient(client *Client) {
	mutex.Lock()
	room.Clients = append(room.Clients, *client)
	client.room = room
	mutex.Unlock()
}

func (room *Room) RemoveClient(remoteIP string) {
	mutex.Lock()
	for i, client := range room.Clients {
		if client.remoteIP == remoteIP {
			room.Clients = append((room.Clients)[:i], (room.Clients)[i+1:]...)
		}
	}
	mutex.Unlock()
}

func CreateRoom(roomList map[string]*Room, roomName string) {
	roomExists := false
	for room := range roomList {
		if room == roomName {
			roomExists = true
			break
		}
	}
	if !roomExists {
		roomList[roomName] = &Room{
			Name:       "\033[36m" + roomName + "\033[0m",
			ParentList: roomList,
			Clients:    []Client{},
			Messages:   []string{},
		}
	}
}

func (client *Client) SwitchRoom(roomName string) {
	if !strings.HasPrefix(roomName, "#") {
		roomName = "#" + roomName
	}

	previousRoom := client.room.Name
	if roomName == toggleRoomName(previousRoom) {
		SendMessage(client.conn, MESSAGE_ERROR_ROOM)
		return
	}

	CreateRoom(client.room.ParentList, roomName)

	client.room.AddMessage(client.username, fmt.Sprintf(MESSAGE_ROOM_LEAVE, toggleRoomName(roomName)))
	client.room.RemoveClient(client.remoteIP)
	client.room.ParentList[roomName].AddClient(client)

	SendMessage(client.conn, MESSAGE_CLEAR)

	for _, message := range client.room.Messages {
		SendMessage(client.conn, message+"\n")
	}

	SendMessage(client.conn, fmt.Sprintf(MESSAGE_CLIENT_SWITCH, client.room.Name, previousRoom))
	client.room.AddMessage(client.username, MESSAGE_ROOM_JOIN)
}
