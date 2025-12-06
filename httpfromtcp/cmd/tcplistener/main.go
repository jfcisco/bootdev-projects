package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	// "strings"
)

const INPUT_BUFFER_SIZE int = 8

var LineBreak byte = '\n'

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)

	// loop goroutine to read lines from `f` 8 bytes at a time
	go func() {
		inputBuffer := make([]byte, INPUT_BUFFER_SIZE)
		var lineBuffer []byte
		for {
			numRead, err := f.Read(inputBuffer)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					log.Fatalf("unexpected error while reading file: %v\n", err)
				}

				// flush line buffer if needed
				if len(lineBuffer) > 0 {
					linesChan <- string(lineBuffer)
				}

				close(linesChan)
				f.Close()
				break
			}

			// scan input buffer for any newlines
			newlineIndex := -1
			for i, inputChar := range inputBuffer[:numRead] {
				if inputChar == LineBreak {
					newlineIndex = i
					break
				}
			}

			if newlineIndex == -1 {
				// no newlines so accumulate everything in line buffer
				lineBuffer = append(lineBuffer, inputBuffer[:numRead]...)
				continue
			}

			// flush line buffer up to newline
			lineBuffer = append(lineBuffer, inputBuffer[0:newlineIndex]...)
			linesChan <- string(lineBuffer)

			// reset line buffer, appending any remaining data in input
			lineBuffer = lineBuffer[:0]
			lineBuffer = append(lineBuffer, inputBuffer[newlineIndex+1:numRead]...)
		}
	}()

	return linesChan
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
			continue
		}

		fmt.Println("connection accepted")

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}

		fmt.Println("connection closed")
	}
}
