package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("Valid single header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "localhost:42069", headers.Get("Host"))
		assert.Equal(t, 23, n)
		assert.False(t, done)
	})

	t.Run("Valid multiple headers", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\nAccept: text/html\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)

		assert.Equal(t, "localhost:42069", headers.Get("Host"))
		assert.Equal(t, 23, n)
		assert.False(t, done)

		n, done, err = headers.Parse(data[n:])
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "text/html", headers.Get("Accept"))
		assert.Equal(t, 19, n)
		assert.False(t, done)
	})

	t.Run("Valid multiple headers with spacing", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n     Accept:    text/html \r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)

		assert.Equal(t, "localhost:42069", headers.Get("Host"))
		assert.Equal(t, 23, n)
		assert.False(t, done)

		n, done, err = headers.Parse(data[n:])
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "text/html", headers.Get("Accept"))
		assert.Equal(t, 28, n)
		assert.False(t, done)
	})

	t.Run("Repeating set headers", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Accept: application/json\r\nAccept: text/html\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.False(t, done)

		n, done, err = headers.Parse(data[n:])
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "application/json, text/html", headers.Get("Accept"))
		assert.Equal(t, 19, n)
		assert.False(t, done)
	})

	t.Run("Repeating set headers - differing casing", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Content-Type: application/json\r\ncontent-type: application/xml\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.False(t, done)

		n, done, err = headers.Parse(data[n:])
		require.NoError(t, err)
		require.NotNil(t, headers)
		assert.Equal(t, "application/json, application/xml", headers.Get("Content-Type"))
		assert.Equal(t, 31, n)
		assert.False(t, done)
	})

	t.Run("Invalid character in field name", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("H©st: localhost:42069\r\n\r\n")
		n, done, err := headers.Parse(data)
		assert.Equal(t, 0, n)
		assert.False(t, done)
		require.Error(t, err)
	})

	t.Run("Valid done", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n\r\n")
		n, _, err := headers.Parse(data)
		require.NoError(t, err)

		n, done, err := headers.Parse(data[n:])
		require.NoError(t, err)
		assert.Equal(t, 0, n)
		assert.True(t, done)
	})

	t.Run("Invalid spacing header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("       Host : localhost:42069       \r\n\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})
}
