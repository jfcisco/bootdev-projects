package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/jfcisco/bootdev-projects/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	open     atomic.Bool
}

func Serve(port int) (*Server, error) {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: ln,
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
	err := response.WriteStatusLine(conn, response.OK)
	if err != nil {
		log.Printf("error writing status line: %s\n", err)
	}

	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		log.Printf("error writing headers: %s\n", err)
	}

	conn.Write([]byte("\r\n"))

	err = conn.Close()
	if err != nil {
		log.Printf("error closing conn: %s\n", err)
	}
}
