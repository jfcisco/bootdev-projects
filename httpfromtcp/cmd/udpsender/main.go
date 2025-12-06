package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalf("error resolving UDP address: %v\n", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("error dialing UDP: %v\n", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		// Show prompt to capture user input
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error while reading line: %s\n", err)
			break
		}

		if _, err := conn.Write([]byte(line)); err != nil {
			fmt.Printf("error writing message: %s\n", err)
		}
	}
}
