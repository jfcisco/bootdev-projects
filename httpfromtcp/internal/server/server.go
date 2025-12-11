package server

import (
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

	w := response.NewWriter()
	s.handler(w, req)

	_, err = w.WriteTo(conn)
	if err != nil {
		log.Printf("error writing headers: %s\n", err)
	}
}
