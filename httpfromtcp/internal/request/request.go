package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/headers"
)

const READER_BUFFER_SIZE = 16
const CRLF = "\r\n"

const (
	Parse_Initialized = iota
	Parse_ParsingHeaders
	Parse_ParsingBody
	Parse_Done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	parseState  int
	contentLen  int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parseSingle(data []byte) (int, error) {
	var readCount int

	switch r.parseState {
	case Parse_Initialized:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		readCount = n
		if reqLine != nil {
			r.RequestLine = *reqLine
			r.parseState = Parse_ParsingHeaders
		}
	case Parse_ParsingHeaders:
		if r.Headers == nil {
			r.Headers = headers.NewHeaders()
		}
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		readCount = n
		if !done {
			return n, nil
		}

		// finished parsing headers
		r.parseState = Parse_ParsingBody
		len, err := strconv.Atoi(r.Headers.Get("Content-Length"))
		if err != nil {
			// for now, if unable to parse content length, assume no body
			len = 0
		}
		r.contentLen = len
	case Parse_ParsingBody:
		remainingLen := r.contentLen - len(r.Body)
		if remainingLen == 0 {
			// finished parsing body
			r.parseState = Parse_Done
			return 0, nil
		}
		if len(data) > remainingLen {
			// content length exceeded
			return 0, errors.New("body length exceeds Content-Length")
		}
		readCount = len(data)
		r.Body = append(r.Body, data[:readCount]...)
	}

	return readCount, nil
}

// Attempts to parse the HTTP request message. If a parse was successful, returns the number of bytes read
func (r *Request) parse(data []byte) (int, error) {
	totalBytes := 0
	for r.parseState != Parse_Done {
		n, err := r.parseSingle(data[totalBytes:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			// not enough data yet
			break
		}

		totalBytes += n
	}
	return totalBytes, nil
}

// Parses the request line from the HTTP message. If successful, returns the RequestLine and number of bytes read from the message
func parseRequestLine(data []byte) (r *RequestLine, readCount int, err error) {
	// Check for presence of CRLF to see if msg is parseable
	crlfIdx := bytes.Index(data, []byte(CRLF))
	if crlfIdx < 0 {
		return nil, 0, nil
	}

	reqLine := string(data[:crlfIdx])

	// request-line should be the first line
	// expect the format (method SP request-target SP method)
	reqLineParts := strings.Split(reqLine, " ")
	if len(reqLineParts) != 3 {
		return nil, 0, errors.New("unexpected number of parts in request-line")
	}

	method := reqLineParts[0]
	target := reqLineParts[1]
	httpVer := reqLineParts[2]

	// validate method is all uppercase
	if method != strings.ToUpper(method) {
		return nil, 0, errors.New("method contains unsupported characters")
	}

	// validate HTTP version and strip prefix "HTTP/"
	if httpVer != "HTTP/1.1" {
		return nil, 0, errors.New("unsupported HTTP version")
	}
	httpVer = httpVer[5:] // we only want the part after HTTP/

	return &RequestLine{
		HttpVersion:   httpVer,
		RequestTarget: target,
		Method:        method,
	}, crlfIdx + len(CRLF), nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	result := &Request{
		parseState: Parse_Initialized,
	}

	var msgBuffer []byte // Holds the request message being parsed
	inputBuf := make([]byte, READER_BUFFER_SIZE)
	eofReached := false // marker to end loop if EOF is reached

	for result.parseState != Parse_Done {
		numRead, err := reader.Read(inputBuf)
		if errors.Is(err, io.EOF) {
			// process remaining data and break at end of current iteration
			eofReached = true
		} else if err != nil {
			return nil, fmt.Errorf("error while reading: %w", err)
		}

		msgBuffer = append(msgBuffer, inputBuf[:numRead]...)

		// Parse HTTP request message
		numParsed, err := result.parse(msgBuffer)
		if err != nil {
			return nil, fmt.Errorf("error while parsing: %w", err)
		}

		if numParsed > 0 {
			msgBuffer = msgBuffer[numParsed:]
		}

		if eofReached {
			break
		}
	}

	if result.parseState != Parse_Done {
		return nil, errors.New("unable to parse HTTP request message")
	}

	return result, nil
}
