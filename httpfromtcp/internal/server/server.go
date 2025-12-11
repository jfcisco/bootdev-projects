package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/request"
	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	handler  Handler
	open     atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: ln,
		handler:  handler,
	}
	s.open.Store(true)

	// start listening
	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.open.Store(false)
	if err := s.listener.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Server) listen() {
	for s.open.Load() {
		conn, err := s.listener.Accept()

		go func() {
			if !s.open.Load() {
				return
			}

			if err != nil {
				log.Printf("error accepting connection: %s\n", err)
				return
			}

			s.handle(conn)
		}()
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("failed to parse request: %s\n", err)
		hErr := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    err.Error(),
		}
		hErr.WriteTo(conn)
		return
	}

	resBuf := new(bytes.Buffer)
	hErr := s.handler(resBuf, req)
	if hErr != nil {
		hErr.WriteTo(conn)
		return
	}

	err = response.WriteStatusLine(conn, response.OK)
	if err != nil {
		log.Printf("error writing status line: %s\n", err)
		return
	}

	err = response.WriteHeaders(conn, response.GetDefaultHeaders(resBuf.Len()))
	if err != nil {
		log.Printf("error writing headers: %s\n", err)
		return
	}

	_, err = resBuf.WriteTo(conn)
	if err != nil {
		log.Printf("error writing body: %s\n", err)
		return
	}
}
