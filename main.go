package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/diegoholiveira/bottle/server"
)

func main() {
	ip := flag.String("address", "0.0.0.0", "IP Address to bind the server")
	port := flag.Int("port", 42000, "Port to bind the server")

	flag.Parse()

	// creates a new server
	server, err := server.New(*ip, *port)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// start handling connections
	go server.Start()

	// Creates a channel to watch for a syscall
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a syscall to stop the server
	<-quit

	// stop the server
	server.Stop()
}
