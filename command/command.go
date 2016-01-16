package command

import (
	"bufio"
	"errors"
	"net"
	"strings"
)

const (
	Use = iota
	Put
	Get
	Purge
	Quit
)

type Command struct {
	Command int
	Data    string
}

func NewCommandFromConnection(conn net.Conn) (*Command, error) {
	scanner := bufio.NewScanner(conn)
	scanner.Scan() // read until the first End-Of-Line

	command := scanner.Text()
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return parseCommand(command)
}

func parseCommand(rawCommand string) (*Command, error) {
	parsed := strings.SplitAfterN(rawCommand, " ", 2)
	parsed[0] = strings.TrimSpace(parsed[0])

	command := new(Command)

	switch parsed[0] {
	case "USE":
		command.Command = Use
		command.Data = parsed[1]
	case "PUT":
		command.Command = Put
		command.Data = parsed[1]
	case "GET":
		command.Command = Get
	case "PURGE":
		command.Command = Purge
	case "QUIT":
		command.Command = Quit
	default:
		return nil, errors.New("Invalid command")
	}

	return command, nil
}
