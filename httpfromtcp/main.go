package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("unable to open messages.txt")
	}

	// Read files 8 bytes at a time
	buf := make([]byte, 8)
	for {
		numRead, err := file.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Fatalf("error while reading file: %v\n", err)
			}
			break
		}
		fmt.Printf("read: %s\n", buf[0:numRead])
	}
}
