package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// io.ReadCloser implementation backed by in-memory buffer
type fakeReadCloser struct {
	buffer bytes.Buffer
}

func NewFakeReadCloser(msg string) *fakeReadCloser {
	fr := &fakeReadCloser{
		buffer: *bytes.NewBufferString(msg),
	}
	return fr
}

func (f *fakeReadCloser) Read(b []byte) (n int, err error) {
	return f.buffer.Read(b)
}

func (f *fakeReadCloser) Close() error {
	return nil
}

func TestGetLinesChannel_Empty(t *testing.T) {
	fr := NewFakeReadCloser("")
	ch := getLinesChannel(fr)

	line := <-ch
	assert.Empty(t, line)
}

func TestGetLinesChannel_OneLine(t *testing.T) {
	fr := NewFakeReadCloser("hello, world!")
	ch := getLinesChannel(fr)

	line := <-ch
	assert.Equal(t, line, "hello, world!")

	_, ok := <-ch
	assert.False(t, ok, "Unexpected data received beyond string")
}

func TestGetLinesChannel_MultiLines(t *testing.T) {
	fr := NewFakeReadCloser(`Do you have what it takes to be an engineer at TheStartup™?
Are you willing to work 80 hours a week in hopes that your 0.001% equity is worth something?
Can you say "synergy" and "democratize" with a straight face?
Are you prepared to eat top ramen at your desk 3 meals a day?
end
`)
	ch := getLinesChannel(fr)

	line := <-ch
	assert.Equal(t, line, "Do you have what it takes to be an engineer at TheStartup™?")

	line = <-ch
	assert.Equal(t, line, "Are you willing to work 80 hours a week in hopes that your 0.001% equity is worth something?")

	line = <-ch
	assert.Equal(t, line, `Can you say "synergy" and "democratize" with a straight face?`)

	line = <-ch
	assert.Equal(t, line, "Are you prepared to eat top ramen at your desk 3 meals a day?")

	line = <-ch
	assert.Equal(t, line, "end")

	_, ok := <-ch
	assert.False(t, ok, "Unexpected data received beyond string")
}

func TestGetLinesChannel_NoTrailingNewline(t *testing.T) {
	fr := NewFakeReadCloser(`{"flavor":"dark mode"}`)
	ch := getLinesChannel(fr)

	line := <-ch
	assert.Equal(t, line, `{"flavor":"dark mode"}`)
}
