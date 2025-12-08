package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse_Good(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestRequestLineParse_GoodWithPath(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestRequestLineParse_GoodPOST(t *testing.T) {
	sr := strings.NewReader("POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/8.17.0\r\nAccept: */*\r\nContent-Type: application/json\r\nContent-Length: 22\r\n\r\n{\"flavor\":\"dark mode\"}")
	r, err := RequestFromReader(sr)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestRequestLineParts_InvalidCount(t *testing.T) {
	// Test: Invalid number of parts in request line
	_, err := RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)
}

func TestRequestLineParse_InvalidMethod(t *testing.T) {
	sr := strings.NewReader("/coffee HTTP/1.1 HEAD")
	_, err := RequestFromReader(sr)
	require.Error(t, err)
}

func TestRequestLineParse_InvalidVersion(t *testing.T) {
	sr := strings.NewReader("PATCH /coffee/1 TEAPOT/1.1")
	_, err := RequestFromReader(sr)
	require.Error(t, err)
}
