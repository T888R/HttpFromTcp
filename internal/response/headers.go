package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	heads := headers.NewHeaders()
	heads.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	heads.Set("Connection", "close")
	heads.Set("Content-Type", "text/plain")
	return heads
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
