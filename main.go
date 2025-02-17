package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	drivers "github.com/siddartha999/F1_Experience_Server/drivers"
)

func main() {
	uniqueDriverEntries := drivers.InitiateDriversInfo()
	uniqueDriverEntryBytes, err := json.Marshal(uniqueDriverEntries)
	if err != nil {
		fmt.Println(err)
	}
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error opening a listening socket", err)
		log.Fatal()
	}
	defer listener.Close()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection")
			break
		}
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n" + string(uniqueDriverEntryBytes) + "\r\n"))
		connection.Close()
	}
}
