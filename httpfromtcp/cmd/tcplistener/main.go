package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/request"
)

func readRequest(rc io.ReadCloser) error {
	defer rc.Close()

	r, err := request.RequestFromReader(rc)
	if err != nil {
		return err
	}

	fmt.Println("Request line:")
	fmt.Printf("- Method: %s\n", r.RequestLine.Method)
	fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

	fmt.Println("connection closed")
	return nil
}

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("unable to listen: %s\n", err)
	}
	defer ln.Close()

	fmt.Println("listening on :42069")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error accepting connection: %s\n", err)
		}

		fmt.Println("connection accepted")

		if err := readRequest(conn); err != nil {
			log.Printf("error from conn: %s\n", err)
			continue
		}
	}
}
