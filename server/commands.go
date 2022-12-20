package server

import (
	"fmt"
	"strings"
)

type Exit bool

type Command struct {
	Name        string
	Description string
	Exec        func(*Client, []string) Exit
}

var commands = map[string]Command{}

var commandHelp = Command{
	Name:        "help",
	Description: "Shows all available commands",
	Exec: func(client *Client, args []string) Exit {
		for _, command := range commands {
			// Print "[color]/[name] - [description]"
			sendMessage(client.conn, fmt.Sprintf("/%v \033[1;30m- %v\033[0m\n", command.Name, command.Description))
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
			client.room.AddMessage("\033[33mServer\033[0m", fmt.Sprintf("\033[33m%v \033[33mis now known as %v", previousUsername, client.username))
		} else {
			sendMessage(client.conn, "Invalid usage: "+strings.TrimSuffix(MESSAGE_USERNAME_ERROR, "Username: "))
		}
		return Exit(false)
	},
}

var commandExit = Command{
	Name:        "exit",
	Description: "Disconnects the user from the server",
	Exec: func(client *Client, args []string) Exit {
		sendMessage(client.conn, MESSAGE_DISCONNECTED)
		pinguSender(client.conn, false)

		client.Disconnect()

		fmt.Fprintln(mw, "User '"+client.username+"' disconnected from the TCP Chat.")
		client.room.AddMessage(client.username, "\033[31mhas left the chat.\033[0m")
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
		sendMessage(client.conn, MESSAGE_HELP)
	}
	return Exit(false)
}

func initCommands() map[string]Command {
	commands = map[string]Command{
		"help": commandHelp,
		"name": commandName,
		"exit": commandExit,
	}
	return commands
}
