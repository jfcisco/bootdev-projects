package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const crlf = "\r\n" // TODO: Move to shared constant

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Get(key string) string {
	if len(key) == 0 {
		return ""
	}
	val, ok := h[strings.ToLower(key)]
	if !ok {
		return ""
	}
	return val
}

func (h Headers) Set(key, value string) {
	currVal := h.Get(key)
	if currVal == "" {
		h[key] = value
		return
	}

	// if a value is already present, append value with comma
	h[key] = fmt.Sprintf("%s, %s", currVal, value)
}

var tokenValidBytes map[byte]struct{}

// Determines whether provided string is a valid token as defined in RFC 9110
func isValidToken(s []byte) bool {
	if tokenValidBytes == nil {
		// initialize set
		tokenValidBytes = make(map[byte]struct{}, 77)

		for i := byte('a'); i <= byte('z'); i++ {
			tokenValidBytes[i] = struct{}{}
		}

		for i := byte('A'); i <= byte('Z'); i++ {
			tokenValidBytes[i] = struct{}{}
		}

		for i := byte('0'); i <= byte('9'); i++ {
			tokenValidBytes[i] = struct{}{}
		}

		// special characters
		tokenValidBytes[byte('!')] = struct{}{}
		tokenValidBytes[byte('#')] = struct{}{}
		tokenValidBytes[byte('$')] = struct{}{}
		tokenValidBytes[byte('%')] = struct{}{}
		tokenValidBytes[byte('&')] = struct{}{}
		tokenValidBytes[byte('\'')] = struct{}{}
		tokenValidBytes[byte('*')] = struct{}{}
		tokenValidBytes[byte('+')] = struct{}{}
		tokenValidBytes[byte('-')] = struct{}{}
		tokenValidBytes[byte('.')] = struct{}{}
		tokenValidBytes[byte('^')] = struct{}{}
		tokenValidBytes[byte('_')] = struct{}{}
		tokenValidBytes[byte('`')] = struct{}{}
		tokenValidBytes[byte('|')] = struct{}{}
		tokenValidBytes[byte('~')] = struct{}{}
	}

	for _, r := range s {
		if _, ok := tokenValidBytes[r]; !ok {
			return false
		}
	}
	return true
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(crlf))
	if crlfIdx < 0 {
		// data not ready for parsing
		return 0, false, nil
	} else if crlfIdx == 0 {
		// leading CRLF indicating end of header section (RFC 9112)
		return len(crlf), true, nil
	}

	// Parse next header
	fLineParts := bytes.Fields(data[:crlfIdx])

	// if this is valid, Fields() should return exactly two parts:
	// [0]: field-name ":"
	// [1]: field-value
	if len(fLineParts) != 2 {
		err = errors.New("parse failed: unexpected number of parts")
		return
	}

	// get field name and preprocess
	fName := fLineParts[0]
	colonIdx := bytes.IndexByte(fName, ':')
	if colonIdx == -1 {
		err = errors.New("parse failed: missing colon in field name")
		return
	}
	fName = fName[:colonIdx]     // re-slice to remove colon
	fName = bytes.ToLower(fName) // store in lowercase since field-names are case-insensitive
	if !isValidToken(fName) {
		return 0, false, errors.New("invalid field name provided")
	}

	// get value from line
	fVal := fLineParts[1]

	// Mutate header
	h.Set(string(fName), string(fVal))

	totalProcessed := crlfIdx + len(crlf)
	return totalProcessed, false, nil
}
