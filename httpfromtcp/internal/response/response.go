package response

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

var reasonPhrases = map[StatusCode]string{
	OK:                  "OK",
	BadRequest:          "Bad Request",
	InternalServerError: "Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var statusLine string
	if phrase, ok := reasonPhrases[statusCode]; ok {
		statusLine = fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, phrase)
	} else {
		statusLine = fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
	}
	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		line := fmt.Appendf(nil, "%s: %s\r\n", key, value)
		_, err := w.Write(line)
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}

type Writer struct {
	status  *bytes.Buffer
	headers *bytes.Buffer
	body    *bytes.Buffer
}

func NewWriter() *Writer {
	return &Writer{
		status:  new(bytes.Buffer),
		headers: new(bytes.Buffer),
		body:    new(bytes.Buffer),
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if err := WriteStatusLine(w.status, statusCode); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if err := WriteHeaders(w.headers, headers); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.body.Write(p)
}

func (w *Writer) WriteTo(dest io.Writer) (int64, error) {
	var totalBytes int64
	n, err := dest.Write(w.status.Bytes())
	if err != nil {
		return 0, err
	}
	totalBytes += int64(n)

	n, err = dest.Write(w.headers.Bytes())
	if err != nil {
		return 0, err
	}
	totalBytes += int64(n)

	n, err = dest.Write(w.body.Bytes())
	if err != nil {
		return 0, err
	}
	totalBytes += int64(n)
	return totalBytes, nil
}
