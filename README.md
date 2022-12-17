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

### Usage

Open your terminal, find the net-cat folder and type `go run .`

This will start the net-cat program on port 8989.

If you want to run this on a different port then type go run . [preferred port]

- Example one `go run .`
  Output: `Listening on :8989`

- Example two `go run . 80`
  Output: `Listening on :80`

To connect to the chat open a new terminal window and type `nc [ip address] [port]`

To exit the chat type `exit`

Following the same examples above.

- Example one `nc 17.6.126.345 8989`

- Example two `nc 17.6.126.345 80`

- Exit the program `exit`

## ⛏️ Built Using <a name = "built_using"></a>

- [Go](https://go.dev/) - Programming language

## ✍️ Authors <a name = "authors"></a>

- [@Falusvampen](https://github.com/Falusvampen)
- [@Kevazy](https://github.com/kevazy)
