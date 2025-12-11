package response

import (
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
