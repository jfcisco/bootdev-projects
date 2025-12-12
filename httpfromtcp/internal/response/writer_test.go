package response

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteChunkedBody(t *testing.T) {
	// Test: Simple chunked body write
	buf := new(bytes.Buffer)
	writer := NewWriter(buf)
	bodyContent := "Hello, World! This is a chunked body test."
	n, err := writer.WriteChunkedBody([]byte(bodyContent))
	require.NoError(t, err)
	require.Greater(t, n, 0)

	assert.Equal(t, "2a\r\nHello, World! This is a chunked body test.\r\n", buf.String())
	buf.Reset() // Clear buffer for next test

	// Test: Multiple chunked body writes
	writer = NewWriter(buf)
	chunk1 := "Hello, "
	chunk2 := "this is "
	chunk3 := "a test."
	n1, err := writer.WriteChunkedBody([]byte(chunk1))
	require.NoError(t, err)
	require.Greater(t, n1, 0)

	n2, err := writer.WriteChunkedBody([]byte(chunk2))
	require.NoError(t, err)
	require.Greater(t, n2, 0)

	n3, err := writer.WriteChunkedBody([]byte(chunk3))
	require.NoError(t, err)
	require.Greater(t, n3, 0)

	expected := "7\r\nHello, \r\n8\r\nthis is \r\n7\r\na test.\r\n"
	assert.Equal(t, expected, buf.String())
	buf.Reset() // Clear buffer for next test

	// Test: Empty chunked body write
	writer = NewWriter(buf)
	n, err = writer.WriteChunkedBody([]byte{})
	require.NoError(t, err)
	require.Equal(t, 5, n)
	assert.Equal(t, "0\r\n\r\n", buf.String())
}
