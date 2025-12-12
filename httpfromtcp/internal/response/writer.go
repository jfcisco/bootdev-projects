package response

import (
	"errors"
	"io"
	"strconv"

	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/headers"
)

type writeState int

const (
	unwritten writeState = iota
	statusWritten
	headersWritten
	bodyWritten
)

type Writer struct {
	conn  io.Writer
	state writeState
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{
		conn:  conn,
		state: unwritten,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != unwritten {
		return errors.New("invalid operation: status line already written")
	}
	if err := WriteStatusLine(w.conn, statusCode); err != nil {
		return err
	}
	w.state = statusWritten
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state > headersWritten {
		return errors.New("invalid operation: headers already written")
	}
	if w.state == unwritten {
		// write a default status line
		if err := WriteStatusLine(w.conn, OK); err != nil {
			return err
		}
	}
	if err := WriteHeaders(w.conn, headers); err != nil {
		return err
	}
	w.state = headersWritten
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state == bodyWritten {
		return 0, errors.New("invalid operation: body already written")
	}
	if w.state < headersWritten {
		// write default headers
		h := GetDefaultHeaders(len(p))
		if err := WriteHeaders(w.conn, h); err != nil {
			return 0, err
		}
	}
	w.state = bodyWritten
	return w.conn.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	var chunk []byte
	// chunk size
	chunk = strconv.AppendInt(chunk, int64(len(p)), 16)
	chunk = append(chunk, crlf...)

	// chunk data
	chunk = append(chunk, p...)
	chunk = append(chunk, crlf...)
	return w.conn.Write(chunk)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.WriteChunkedBody([]byte{})
}
