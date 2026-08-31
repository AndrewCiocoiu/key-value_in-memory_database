package main

import (
	"log"
	"net"
)

const PORT = ":9000"

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	defer listener.Close()

	log.Printf("Server is listening on port: %s\n", PORT[1:])

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		connection.Write([]byte("Hello, world!\n"))

		connection.Close()
	}
}
