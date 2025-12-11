package server

import (
	"io"

	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/request"
	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (e *HandlerError) WriteTo(w io.Writer) (int64, error) {
	if err := response.WriteStatusLine(w, e.StatusCode); err != nil {
		return 0, err
	}

	msgBytes := []byte(e.Message)
	if err := response.WriteHeaders(w, response.GetDefaultHeaders(len(msgBytes))); err != nil {
		return 0, err
	}

	if _, err := w.Write(msgBytes); err != nil {
		return 0, err
	}

	return int64(len(msgBytes)), nil
}
