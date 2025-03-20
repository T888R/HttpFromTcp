package main

import (
	"errors"
	// "bufio"
	"fmt"
	"io"
	"log"
	"net"
	// "time"
	// "os"
	"strings"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("Could not create listener", err)
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", connection.RemoteAddr())

		msgChan := getLinesChannel(connection)

		for line := range msgChan {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", connection.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	lines := make(chan string)
//
// 	go func() {
// 		defer f.Close()
// 		defer close(lines)
//
// 		// Use a scanner with a deadline to avoid hanging
// 		scanner := bufio.NewScanner(f)
//
// 		// If f is a net.Conn, set a read deadline
// 		if conn, ok := f.(net.Conn); ok {
// 			conn.SetReadDeadline(time.Now().Add(2 * time.Second))
// 		}
//
// 		for scanner.Scan() {
// 			lines <- scanner.Text()
// 		}
//
// 		// Check if there's any content left in the buffer
// 		if lineContents := scanner.Text(); lineContents != "" {
// 			lines <- lineContents
// 		}
//
// 		if err := scanner.Err(); err != nil && err != io.EOF {
// 			fmt.Printf("error: %s\n", err.Error())
// 		}
// 	}()
//
// 	return lines
// }
