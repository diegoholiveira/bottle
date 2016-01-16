package server

import (
	"fmt"
	"net"
	"os"

	"github.com/diegoholiveira/bottle/command"
)

type Server struct {
	address *net.TCPAddr
}

func (server *Server) Init(address string, port int) {
	ip := net.ParseIP(address)
	if ip == nil {
		fmt.Println("You must define a valid IP for bottle")
		os.Exit(1)
	}

	server.address = &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
}

func (server *Server) Start() {
	fmt.Printf("Starting bottle at %s\n", server.address.String())

	listener, err := net.ListenTCP("tcp", server.address)
	if err != nil {
		fmt.Println("Its not possible to start bottle, please verify", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	for {
		comm, err := command.NewCommandFromConnection(conn)
		if err != nil {
			conn.Write([]byte(err.Error()))
			return
		}

		var msg []byte

		switch comm.Command {
		default:
			return
		case command.Put:
			msg = []byte("Putting an item...\n")
		case command.Get:
			msg = []byte("Getting an item...\n")
		case command.Use:
			msg = []byte("Using a queue...\n")
		case command.Purge:
			msg = []byte("Purging a queue\n")
		}
		conn.Write(msg)
	}
}
