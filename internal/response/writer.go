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
