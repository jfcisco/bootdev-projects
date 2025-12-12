package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/headers"
	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/request"
	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/response"
)

func routeHandler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		httpBinProxyHandler(w, req)
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		videoHandler(w)
	} else {
		htmlHandler(w, req)
	}
}

func writeHtmlContent(w *response.Writer, html string) {
	h := response.GetDefaultHeaders(len(html))
	h.SetNoAppend("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(html))
}

func htmlHandler(w *response.Writer, req *request.Request) {
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

// var httpBinBaseURL = "http://httpbin.org"
var httpBinBaseURL = "http://localhost:8081"

func httpBinProxyHandler(w *response.Writer, req *request.Request) {
	// Dynamic routes
	after, ok := strings.CutPrefix(req.RequestLine.RequestTarget, "/httpbin")
	if !ok {
		w.WriteStatusLine(response.BadRequest)
		w.WriteBody([]byte("no count provided"))
		return
	}

	log.Printf("received /httpbin/%s\n", after)
	w.WriteStatusLine(response.OK)

	h := headers.NewHeaders()
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Connection", "close")
	h.Set("Trailer", "x-content-sha256")
	h.Set("Trailer", "x-content-length")
	w.WriteHeaders(h)

	binResp, err := getStreamFromHttpBin(after)
	if err != nil {
		log.Printf("error fetching httpbin stream: %v\n", err)
		w.WriteStatusLine(response.InternalServerError)
		w.WriteBody([]byte("error fetching httpbin stream"))
		return
	}
	defer binResp.Body.Close()

	log.Println("writing response")
	hasher := sha256.New()
	rawBodyLength := 0
	buf := make([]byte, 1024)
	for {
		n, err := binResp.Body.Read(buf)
		if n > 0 {
			_, err := w.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Printf("Error writing chunked body: %v\n", err)
				break
			}
			rawBodyLength += n
			hasher.Write(buf[:n])
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

	// write trailers
	trailers := headers.NewHeaders()
	contentHash := hasher.Sum(nil)
	trailers.Set("X-Content-Sha256", hex.EncodeToString(contentHash))
	trailers.Set("X-Content-Length", strconv.Itoa(rawBodyLength))
	w.WriteTrailers(trailers)
	log.Println("done writing trailers")
}

func getStreamFromHttpBin(endpoint string) (*http.Response, error) {
	return http.Get(httpBinBaseURL + endpoint)
}

func videoHandler(w *response.Writer) {
	if err := w.WriteStatusLine(response.OK); err != nil {
		w.WriteStatusLine(response.InternalServerError)
		w.WriteBody([]byte("error writing status line"))
		return
	}

	h := headers.NewHeaders()
	h.Set("Content-Type", "video/mp4")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Connection", "close")
	if err := w.WriteHeaders(h); err != nil {
		w.WriteStatusLine(response.InternalServerError)
		w.WriteBody([]byte("error writing headers"))
		return
	}

	// read video file into memory
	f, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		w.WriteStatusLine(response.InternalServerError)
		w.WriteBody([]byte("error reading video file"))
		return
	}
	if _, err := w.WriteChunkedBody(f); err != nil {
		w.WriteStatusLine(response.InternalServerError)
		w.WriteBody([]byte("error writing chunked body"))
		return
	}
	if _, err := w.WriteChunkedBodyDone(); err != nil {
		w.WriteStatusLine(response.InternalServerError)
		w.WriteBody([]byte("error writing chunked body done"))
		return
	}
	if err := w.WriteTrailers(nil); err != nil {
		w.WriteStatusLine(response.InternalServerError)
		w.WriteBody([]byte("error writing trailers"))
		return
	}
}
