package headers

import (
	"bytes"
	"errors"
)

const crlf = "\r\n" // TODO: Move to shared constant

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(crlf))
	if crlfIdx < 0 {
		// data not ready for parsing
		return 0, false, nil
	} else if crlfIdx == 0 {
		// leading CRLF indicating end of header section (RFC 9112)
		return 0, true, nil
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

	fName := fLineParts[0]
	colonIdx := bytes.IndexByte(fName, ':')
	if colonIdx == -1 {
		err = errors.New("parse failed: missing colon in field name")
		return
	}
	fName = fName[:colonIdx] // re-slice to remove colon
	fVal := fLineParts[1]

	// Mutate header
	// TODO: handle difference in casing here?
	h[string(fName)] = string(fVal)

	totalProcessed := crlfIdx + len(crlf)
	return totalProcessed, false, nil
}
