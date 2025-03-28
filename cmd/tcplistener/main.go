package main

import (
	"httpfromtcp/internal/request"
	// "bufio"
	"fmt"
	// "io"
	"log"
	"net"
	// "time"
	// "os"
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

		req, err := request.RequestFromReader(connection)
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Printf(
			"Request line:\n - Method: %s\n - Target: %s\n - Version: %s\n",
			req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion,
		)

		fmt.Println("Connection to ", connection.RemoteAddr(), "closed")
	}
}
