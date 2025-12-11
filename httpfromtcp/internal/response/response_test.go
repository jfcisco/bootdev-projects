package response

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteStatusLine(t *testing.T) {
	// Test: 200 OK
	buf := new(bytes.Buffer)
	WriteStatusLine(buf, OK)
	actual, _ := io.ReadAll(buf)
	assert.Equal(t, "HTTP/1.1 200 OK\r\n", string(actual))

	// Test: 400 Bad Request
	buf.Reset()
	WriteStatusLine(buf, BadRequest)
	actual, _ = io.ReadAll(buf)
	assert.Equal(t, "HTTP/1.1 400 Bad Request\r\n", string(actual))

	// Test: Non-standard status code
	buf.Reset()
	WriteStatusLine(buf, 999)
	actual, _ = io.ReadAll(buf)
	assert.Equal(t, "HTTP/1.1 999 \r\n", string(actual))
}
