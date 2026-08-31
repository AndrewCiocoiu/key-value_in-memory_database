package main

import (
	"bufio"
	"fmt"
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
		fmt.Fprintf(connection, "Welcome to Andrew's in-memory DB!\n")

		reader := bufio.NewReader(connection)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}

			fmt.Fprintf(connection, "%s", line)
		}
		connection.Close()
	}
}
