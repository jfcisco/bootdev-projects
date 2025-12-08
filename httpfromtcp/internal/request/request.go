package request

import (
	"errors"
	"io"
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

const CRLF string = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	msg, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	msgLines := strings.Split(string(msg), CRLF)

	if len(msgLines) == 0 {
		return nil, errors.New("read empty request")
	}

	// request-line should be the first line
	// expect the format (method SP request-target SP method)
	reqLine := msgLines[0]
	reqLineParts := strings.Split(reqLine, " ")

	if len(reqLineParts) != 3 {
		return nil, errors.New("unexpected number of parts in request-line")
	}

	method := reqLineParts[0]
	target := reqLineParts[1]
	httpVer := reqLineParts[2]

	// validate method is all uppercase
	if method != strings.ToUpper(method) {
		return nil, errors.New("method contains unsupported characters")
	}

	// validate HTTP version and strip prefix "HTTP/"
	if httpVer != "HTTP/1.1" {
		return nil, errors.New("unsupported HTTP version")
	}
	httpVer = httpVer[5:] // we only want the part after HTTP/

	result := &Request{
		RequestLine: RequestLine{
			HttpVersion:   httpVer,
			RequestTarget: target,
			Method:        method,
		},
	}
	return result, nil
}
