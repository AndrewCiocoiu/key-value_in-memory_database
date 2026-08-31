package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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

		//Strip newline and carriage return on windows
		clean_string := strings.TrimSpace(line)

		words := strings.Split(clean_string, " ")

		if len(words) > 3 {
			fmt.Fprintf(conn, "ERROR: Bad format.\n")
			continue
		}

		switch command := words[0]; command {
		case "SET":
			if len(words) != 3 {
				fmt.Fprintf(conn, "ERROR: Bad format.\n")
				continue
			}

			mu.Lock()
			db[words[1]] = words[2]
			fmt.Fprintf(conn, "Value set!\n")
			mu.Unlock()
		case "GET":
			if len(words) != 2 {
				fmt.Fprintf(conn, "ERROR: Bad format.\n")
				continue
			}

			mu.RLock()

			if val, ok := db[words[1]]; ok {
				fmt.Fprintf(conn, "%s\n", val)
			} else {
				fmt.Fprintf(conn, "Key does not exist in DB!\n")
			}

			mu.RUnlock()
		case "DEL":
			if len(words) != 2 {
				fmt.Fprintf(conn, "ERROR: Bad format.\n")
				continue
			}

			mu.Lock()
			if _, ok := db[words[1]]; ok {
				delete(db, words[1])
				fmt.Fprintf(conn, "%s Sucesfully deleted!\n", words[1])
			} else {
				fmt.Fprintf(conn, "Key does not exist in DB!\n")
			}
			mu.Unlock()
		case "SAVE":
			if len(words) != 1 {
				fmt.Fprintf(conn, "ERROR: Bad format.\n")
				continue
			}

			file, err := os.Create("dump.txt")
			if err != nil {
				fmt.Fprintf(conn, "There was an error saving!\n%s\n", err)
				continue
			}

			mu.RLock()

			for key, val := range db {
				fmt.Fprintf(file, "%s, %s\n", key, val)
			}
			file.Close()

			fmt.Fprintf(conn, "Database dumped in dump.txt succesfully!\n")

			mu.RUnlock()
		default:
			fmt.Fprintf(conn, "ERROR: Unknown command.\n")
		}

	}

}

func loadDBFromDump() {
	file, err := os.Open("dump.txt")
	if err != nil {
		return
	}

	reader := bufio.NewScanner(file)

	for reader.Scan() {
		line := reader.Text()
		entry := strings.Split(line, ", ")
		if len(entry) == 2 {
			db[entry[0]] = entry[1]
		}
	}

	if err := reader.Err(); err != nil {
		log.Printf("There was an error restoring the DB. DB will not be restored.\n")
		db = make(map[string]string)
		return
	}

	log.Printf("Loaded Database from disk!\n")

}

func main() {

	loadDBFromDump()

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
