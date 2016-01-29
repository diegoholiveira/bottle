package server

import (
	"fmt"
	"net"
	"os"
	"unicode/utf8"

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

		go handle(server, conn)
	}
}

func handle(server *Server, conn net.Conn) {
	defer conn.Close()

	var q *queue.Queue

	for {
		comm, err := command.NewCommandFromConnection(conn)
		if err != nil {
			conn.Write([]byte(err.Error()))
			break
		}

		var msg []byte

		if comm.Command == command.Quit {
			break
		}

		if comm.Command != command.Use && q == nil {
			conn.Write([]byte("Select a queue first"))
			break
		}

		msg = []byte("OK")

		switch comm.Command {
		case command.Put:
			q.Lock()
			q.Push(comm.Data)
			q.Unlock()
		case command.Get:
			q.Lock()
			msg = []byte("NULL")
			if q.Len() > 0 {
				if item := q.Pop(); utf8.RuneCountInString(item) > 0 {
					msg = []byte(item)

				}
			}
			q.Unlock()
		case command.Use:
			if _, ok := server.queues[comm.Data]; !ok {
				fmt.Printf("Creating a queue named %s\n", comm.Data)

				server.queues[comm.Data] = queue.New()
			}

			q = server.queues[comm.Data]
		case command.Purge:
			// TODO: implement it in a way that we can lock it!
			q = queue.New()
		}
		conn.Write(msg)
	}
}
