package comms

import (
	"fmt"
	"log"
	"net"
)

type TCP struct {
	socket
	listener net.Listener
	conn     net.Conn
}

func NewTCP(ip string, port string) *TCP {
	return &TCP{
		socket: socket{
			ip:   ip,
			port: port,
			io:   nil,
		},
		listener: nil,
		conn:     nil,
	}
}

func (s *TCP) Close() {
	if s.listener != nil {
		s.listener.Close()
	}
	s.conn.Close()
}

func (s *TCP) Server() {
	var err error
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("Error while listening tcp connection: %v\n", err)
	}

	s.conn, err = s.listener.Accept()
	s.io = s.conn
	if err != nil {
		log.Panicf("Error while accepting tcp connection: %v\n", err)
	}
}

func (s *TCP) Client() {
	var err error
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	s.conn, err = net.Dial("tcp", addr)
	s.io = s.conn
	if err != nil {
		log.Panicf("Error while dialing tcp connection: %v\n", err)
	}
}
