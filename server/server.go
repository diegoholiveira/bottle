package server

import (
	"errors"
	"log"
	"net"
	"os"
	"sync"
	"unicode/utf8"

	"github.com/diegoholiveira/bottle/command"
	"github.com/diegoholiveira/bottle/queue"
)

type queues map[string]*queue.Queue

type Server struct {
	address   *net.TCPAddr
	waitGroup *sync.WaitGroup
	queues    queues
	listener  *net.TCPListener
}

func New(address string, port int) (*Server, error) {
	ip := net.ParseIP(address)
	if ip == nil {
		return nil, errors.New("You must define a valid IP for Bottle")
	}

	server := &Server{
		waitGroup: &sync.WaitGroup{},
		address: &net.TCPAddr{
			IP:   ip,
			Port: port,
		},
		queues: make(queues),
	}

	return server, nil
}

func (server *Server) Start() {
	log.Printf("Starting bottle at %s\n", server.address.String())

	listener, err := net.ListenTCP("tcp", server.address)
	if err != nil {
		log.Println(err.Error())
		log.Println("It's not possible to start bottle!")
		os.Exit(1)
	}

	server.listener = listener

	for {
		// keep accepting new connections
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		log.Println(conn.RemoteAddr(), "connected")

		// adds the new client to the wait group
		server.waitGroup.Add(1)
		// handle with the new client connection
		go server.handle(conn)
	}
}

func (server *Server) Stop() {
	log.Println("Stopping the server...")
	server.listener.Close()
	server.waitGroup.Wait()
}

func (server *Server) handle(conn net.Conn) {
	defer conn.Close()
	defer server.waitGroup.Done()

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
				log.Printf("Creating a queue named %s\n", comm.Data)

				server.queues[comm.Data] = queue.New()
			}

			q = server.queues[comm.Data]
		case command.Purge:
			// TODO: implement it in a way that we can lock it!
			q = queue.New()
		}
		conn.Write(msg)
	}

	log.Println(conn.RemoteAddr(), "is done")
}
