package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "142.132.167.112:25565")
	if err != nil {
		log.Fatalf("Failed to resolve TCP address: %s\n", err)
	}
	con, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to dial TCP connection: %s\n", err)
	}

	err = con.Close()
	if err != nil {
		log.Fatalf("Failed to close connection: %s\n", err)
	}

	fmt.Println("Works")
}
