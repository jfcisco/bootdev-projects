package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var LineBreak byte = '\n'

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("unable to open messages.txt")
	}

	// to read data 8 bytes at a time
	inputBuffer := make([]byte, 8)
	var lineBuffer []byte
	for {
		numRead, err := file.Read(inputBuffer)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Fatalf("error while reading file: %v\n", err)
			}
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

		// newline found so flush buffer to stdout
		lineBuffer = append(lineBuffer, inputBuffer[0:newlineIndex]...)
		fmt.Printf("read: %s\n", lineBuffer)
		lineBuffer = lineBuffer[:0]
		lineBuffer = append(lineBuffer, inputBuffer[newlineIndex+1:numRead]...)
	}
}
