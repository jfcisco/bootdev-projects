package main

import (
	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/request"
	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/response"
)

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
