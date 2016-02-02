package server

import (
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/diegoholiveira/bottle/command"
	"github.com/diegoholiveira/bottle/queue"
)

type queues map[string]*queue.Queue

type Server struct {
	address   *net.TCPAddr
	waitGroup *sync.WaitGroup
	queues    queues
	quit      chan bool
}

func New(address string, port int) (*Server, error) {
	p := strconv.Itoa(port)

	ip, err := net.ResolveTCPAddr("tcp4", address+":"+p)
	if err != nil {
		return nil, err
	}

	server := &Server{
		waitGroup: &sync.WaitGroup{},
		address:   ip,
		queues:    make(queues),
		quit:      make(chan bool),
	}

	return server, nil
}

func (server *Server) Start() {
	log.Printf("Starting bottle at %s\n", server.address.String())

	listener, err := net.ListenTCP("tcp4", server.address)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	defer listener.Close()

	// Define a timeout for listen new connections
	listener.SetDeadline(time.Now().Add(5 * time.Second))

	for {
		select {
		case <-server.quit:
			return
		default:
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
}

func (server *Server) Stop() {
	log.Println("Stopping the server...")

	// Stop listener for new connections
	server.quit <- true

	// Blocks the execution until all clients disconnect
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
