package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var db = make(map[string]string)
var mu sync.RWMutex

const PORT = ":9000"

func handleConnect(conn net.Conn, mu sync.RWMutex) {
	defer conn.Close()

	fmt.Fprintf(conn, "Welcome to Andrew's in-memory DB!\n")

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		clean_string := line[:len(line)-3]

		words := strings.Split(clean_string, " ")

		switch command := words[0]; command {
		case "SET":
		case "GET":
		case "DEL":
		default:
			fmt.Fprintf(conn, "ERROR: Unknown command.\n")
		}

	}

}

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
			log.Printf("%s\n", err)
			continue
		}

		go handleConnect(connection, mu)

	}
}
