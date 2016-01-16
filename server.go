package main

import (
	"fmt"
	"net"
	"os"
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

		server.Handle(conn)
	}
}

func (server *Server) Handle(conn net.Conn) {
	defer conn.Close()

	msg := []byte("Hello from bottle\n")
	conn.Write(msg)
}
