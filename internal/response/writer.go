package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type writerState int

const (
	writeStatusLine writerState = iota
	writeStatusHeaders
	writeStatusBody
	writeStatusTrailers
)

const (
	crlf       = "\r\n"
	bufferSize = 8
)

type Writer struct {
	writerState writerState
	writer      io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: writeStatusLine,
		writer:      w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writeStatusLine {
		return fmt.Errorf("Can't write status line %d", w.writerState)
	}
	defer func() { w.writerState = writeStatusHeaders }()
	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != writeStatusHeaders {
		return fmt.Errorf("Can't write headers %d", w.writerState)
	}
	defer func() { w.writerState = writeStatusBody }()
	for k, v := range h {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writeStatusBody {
		return 0, fmt.Errorf("Can't write body %d", w.writerState)
	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writeStatusBody {
		return 0, fmt.Errorf("Can't write chunked body %d", w.writerState)
	}
	chunkSize := len(p)
	nTotal := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != writeStatusBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}
	n, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}
	w.writerState = writeStatusTrailers
	return 0, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != writeStatusTrailers {
		return fmt.Errorf("Can't write trailers in state %d", w.writerState)
	}
	defer func() { w.writerState = writeStatusBody }()
	for k, v := range h {
		fmt.Printf("%s: %s\r\n", k, v)
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
