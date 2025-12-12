package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/headers"
	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/request"
	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/response"
)

func routeHandler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream/") {
		handleStream(w, req)
		return
	}
	handle(w, req)
}

func writeHtmlContent(w *response.Writer, html string) {
	h := response.GetDefaultHeaders(len(html))
	h.SetNoAppend("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(html))
}

func handle(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.BadRequest)
		content := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
		writeHtmlContent(w, content)
	case "/myproblem":
		w.WriteStatusLine(response.InternalServerError)
		content := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`
		writeHtmlContent(w, content)
	default:
		w.WriteStatusLine(response.OK)
		content := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`
		writeHtmlContent(w, content)
	}
}

var httpBinBaseURL = "http://httpbin.org"

func handleStream(w *response.Writer, req *request.Request) {
	// Dynamic routes
	after, ok := strings.CutPrefix(req.RequestLine.RequestTarget, "/httpbin/stream/")
	if !ok {
		w.WriteStatusLine(response.BadRequest)
		w.WriteBody([]byte("no count provided"))
		return
	}

	count, err := strconv.Atoi(after)
	if err != nil {
		w.WriteStatusLine(response.BadRequest)
		w.WriteBody([]byte("invalid count provided"))
		return
	}

	log.Printf("received /httpbin/%d\n", count)
	w.WriteStatusLine(response.OK)

	h := headers.NewHeaders()
	h.Set("Content-Type", "application/json")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Connection", "close")
	w.WriteHeaders(h)

	streamUrl := httpBinBaseURL + "/stream/" + strconv.Itoa(count)
	log.Println("fetching stream from", streamUrl)
	httpResp, err := http.Get(streamUrl)
	if err != nil {
		log.Printf("Error fetching httpbin stream: %v\n", err)
		w.WriteStatusLine(response.InternalServerError)
		w.WriteBody([]byte("error fetching httpbin stream"))
		return
	}
	defer httpResp.Body.Close()

	log.Println("writing response")
	buf := make([]byte, 1024)
	for {
		n, err := httpResp.Body.Read(buf)
		if n > 0 {
			_, err := w.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Printf("Error writing chunked body: %v\n", err)
				break
			}
		}
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("Error writing chunk: %v\n", err)
			}
			break
		}
	}
	w.WriteChunkedBodyDone()
	log.Println("done writing chunked response")
}
