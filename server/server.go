package server

import (
	"fmt"
	"net"
	"os"

	"github.com/diegoholiveira/bottle/command"
	"github.com/diegoholiveira/bottle/queue"
)

type queues map[string]*queue.Queue

type Server struct {
	address *net.TCPAddr
	queues  queues
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

	server.queues = make(queues)
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

		server.handle(conn)
	}
}

func (server *Server) handle(conn net.Conn) {
	defer conn.Close()

	var q *queue.Queue

	for {
		comm, err := command.NewCommandFromConnection(conn)
		if err != nil {
			conn.Write([]byte(err.Error()))
			return
		}

		var msg []byte

		if comm.Command == command.Quit {
			return
		}

		if comm.Command != command.Use && q == nil {
			conn.Write([]byte("Select a queue first"))

			continue
		}

		msg = []byte("OK")

		switch comm.Command {
		case command.Put:
			q.Push(comm.Data)
		case command.Get:
			if q.Len() == 0 {
				msg = []byte("NULL")
			} else {
				item := q.Pop()
				msg = []byte(item)
			}
		case command.Use:
			if _, ok := server.queues[comm.Data]; !ok {
				fmt.Printf("Creating a queue named %s\n", comm.Data)

				server.queues[comm.Data] = queue.New()
			}

			q = server.queues[comm.Data]
		case command.Purge:
			q = queue.New()
		}
		conn.Write(msg)
	}
}
