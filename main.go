package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Could not open message file", err)
	}

	linesChan := getLinesChannel(file)

	for line := range linesChan {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)
		lineContents := ""
		for {
			buf := make([]byte, 8, 8)
			n, err := f.Read(buf)
			if err != nil {
				// fmt.Println("Error reading:", err)
				if lineContents != "" {
					// fmt.Printf("read: %s\n", lineContents)
					lines <- lineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			str := string(buf[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", lineContents, parts[i])
				lineContents = ""
			}
			lineContents += parts[len(parts)-1]
		}
	}()

	return lines
}
