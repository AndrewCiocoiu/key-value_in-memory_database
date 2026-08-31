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

func handleConnect(conn net.Conn, mu *sync.RWMutex) {
	defer conn.Close()

	fmt.Fprintf(conn, "Welcome to Andrew's in-memory DB!\n")

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		clean_string := line[:len(line)-1]

		words := strings.Split(clean_string, " ")

		if len(words) > 3 {
			fmt.Fprintf(conn, "ERROR: Bad format.\n")
			continue
		}

		switch command := words[0]; command {
		case "SET":
			mu.Lock()
			db[words[1]] = words[2]
			fmt.Fprintf(conn, "Value set!\n")
			mu.Unlock()
		case "GET":
			mu.RLock()

			if val, ok := db[words[1]]; ok {
				fmt.Fprintf(conn, "%s\n", val)
			} else {
				fmt.Fprintf(conn, "Key does not exist in DB!\n")
			}

			mu.RUnlock()
		case "DEL":
			mu.Lock()
			if _, ok := db[words[1]]; ok {
				delete(db, words[1])
				fmt.Fprintf(conn, "%s Sucesfully deleted!\n", words[1])
			} else {
				fmt.Fprintf(conn, "Key does not exist in DB!\n")
			}
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

		go handleConnect(connection, &mu)

	}
}
