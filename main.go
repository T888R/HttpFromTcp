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
		os.Exit(1)
	}

	lineContents := ""

	for {
		buf := make([]byte, 8, 8)
		n, err := file.Read(buf)
		if err != nil {
			// fmt.Println("Error reading:", err)
			if lineContents != "" {
				fmt.Printf("read: %s\n", lineContents)
				lineContents = ""
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
			fmt.Printf("read: %s%s\n", lineContents, parts[i])
			lineContents = ""
		}
		lineContents += parts[len(parts)-1]
		// fmt.Println("read:", str)
	}
}
