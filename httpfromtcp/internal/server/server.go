package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
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
	msg := "HTTP/1.1 200 OK\r\n" + // status line
		"Content-Type: text/plain\r\n" + // headers start
		"Content-Length: 13\r\n" +
		"\r\n" + // empty line to indicate headers end
		"Hello World!\n" // body
	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Printf("error writing response: %s\n", err)
	}

	err = conn.Close()
	if err != nil {
		log.Printf("error closing conn: %s\n", err)
	}
}
