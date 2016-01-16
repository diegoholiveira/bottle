package main

import (
	"flag"
)

func main() {
	ip := flag.String("address", "0.0.0.0", "IP Address to bind the server")
	port := flag.Int("port", 42000, "Port to bind the server")

	flag.Parse()

	server := new(Server)
	server.Init(*ip, *port)
	server.Start()
}
