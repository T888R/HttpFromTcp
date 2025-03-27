package request

import (
	"errors"
	// "fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(r io.Reader) (*Request, error) {

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	parts := strings.Split(string(b), "\r\n")
	fields := strings.Fields(parts[0])

	if len(fields) != 3 {
		return nil, errors.New("Incorrect number of fields in request line")
	}

	rawHttpVersion := fields[2]
	splitVersion := strings.Split(rawHttpVersion, "HTTP/")
	httpVersion := splitVersion[1]
	reqTarget := fields[1]
	reqMethod := strings.ToUpper(fields[0])

	if splitVersion[1] != "1.1" {
		return nil, errors.New("Http version is not supported")
	}

	reqLine := RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: reqTarget,
		Method:        reqMethod,
	}

	req := Request{RequestLine: reqLine}

	return &req, nil
}
