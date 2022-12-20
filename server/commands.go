package server

import (
	"fmt"
	"strings"
)

type Exit bool

type Command struct {
	Name        string
	Description string
	Exec        func(client *Client, args []string) Exit
}

var commands = map[string]Command{}

var commandHelp = Command{
	Name:        "help",
	Description: "Shows all available commands",
	Exec: func(client *Client, args []string) Exit {
		for _, command := range commands {
			// Print "[color]/[name] - [description]"
			SendMessage(client.conn, fmt.Sprintf(MESSAGE_COMMAND_HELP, command.Name, command.Description))
		}
		return Exit(false)
	},
}

var commandRoom = Command{
	Name:        "room",
	Description: "Switches the user's chatroom: /room <room_name>",
	Exec: func(client *Client, args []string) Exit {
		if len(args) == 1 {
			client.SwitchRoom(args[0])
		} else {
			SendMessage(client.conn, MESSAGE_ERROR_USAGE_ROOM)
		}
		return Exit(false)
	},
}

var commandRooms = Command{
	Name:        "rooms",
	Description: "Lists all rooms",
	Exec: func(client *Client, args []string) Exit {
		for _, room := range client.room.ParentList {
			SendMessage(client.conn, fmt.Sprintf(MESSAGE_COMMAND_ROOMS, room.Name, len(room.Clients)))
		}
		return Exit(false)
	},
}

var commandName = Command{
	Name:        "name",
	Description: "Change username: /name <new_name>",
	Exec: func(client *Client, args []string) Exit {
		previousUsername := client.username
		newUsername := strings.Join(args, " ")
		if len(args) > 0 && validNameBool(newUsername) {
			client.username = randomizeColor() + newUsername + "\033[0m"
			client.room.AddMessage(previousUsername, fmt.Sprintf(MESSAGE_ACTION_NAME, client.username))
		} else {
			SendMessage(client.conn, "Invalid usage: "+strings.TrimSuffix(MESSAGE_ERROR_USERNAME, "Username: "))
		}
		return Exit(false)
	},
}

var commandExit = Command{
	Name:        "exit",
	Description: "Disconnects the user from the server",
	Exec: func(client *Client, args []string) Exit {
		SendMessage(client.conn, MESSAGE_CLIENT_DISCONNECTED)
		pinguSender(client.conn, false)

		client.Disconnect()

		fmt.Fprintln(mw, "User '"+client.username+"' disconnected from the TCP Chat.")
		client.room.AddMessage(client.username, MESSAGE_ACTION_LEAVE)
		return Exit(true)
	},
}

func commandHandler(client *Client, input string) Exit {
	// Initialize commands here to avoid invalid initialization
	initCommands()
	args := strings.Split(input, " ")
	if command, ok := commands[args[0]]; ok {
		fmt.Fprintf(mw, "[%v][%v] executed commmand: /%v\n", getCurrentTime(), client.username, strings.Join(args, " "))
		return command.Exec(client, args[1:])
	} else {
		SendMessage(client.conn, MESSAGE_ERROR_HELP)
	}
	return Exit(false)
}

func initCommands() map[string]Command {
	commands = map[string]Command{
		"help":  commandHelp,
		"room":  commandRoom,
		"rooms": commandRooms,
		"name":  commandName,
		"exit":  commandExit,
	}
	return commands
}
