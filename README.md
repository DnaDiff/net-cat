
# net-cat

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Built Using](#⛏️-built-using)
- [Authors](#✍️-authors)

## About <a name = "about"></a>

The purpose of this project is to develop a functional netcat program to enable communication over a network as part of the grit:lab coursework.
This is the fifth project in the series.
## Getting Started <a name = "getting_started"></a>

Download the repository to your local machine.

### Prerequisites

- [Go](https://go.dev/) 1.19
- A **Unix-based** terminal

### Usage

#### Server
Open your terminal, find the net-cat folder and run the command `go run .`

This will start the net-cat program on port 8989.

If you want to run this on a different port then run `go run . [preferred port]` instead.

- Example one `go run .`
  Output: `Listening on :8989`

- Example two `go run . 80`
  Output: `Listening on :80`

You can also enable logging of all messages and server output by including the flag `--log`

- Example one `go run . 8080 --log`
- Example two `go run . --log`

These logs should appear inside a `/logs/` folder inside the project location.

#### Client
Connect to the chat by opening a new **Unix-based** terminal window and run `nc [ip address] [port]`

To view available commands, type `/help`

If you'd like to change your username, you can type `/name <new_name>`

Join different chatrooms by typing `/room <room_name>`

To exit the chat type `/exit`

Following the same examples above.

- Example one `nc 17.6.126.345 8989`

- Example two `nc 17.6.126.345 80`

- Change room `/room #new-room`

- Exit the program `/exit`

## ⛏️ Built Using <a name = "built_using"></a>

- [Go](https://go.dev/) - Programming language

## ✍️ Authors <a name = "authors"></a>

- [@Falusvampen](https://github.com/Falusvampen)
- [@Kevazy](https://github.com/kevazy)
