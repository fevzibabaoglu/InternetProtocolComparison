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

func newTCP(ip string, port string) *TCP {
	return &TCP{
		socket: socket{
			ip:    ip,
			port:  port,
			io:    nil,
			chans: newChannels(),
		},
		listener: nil,
		conn:     nil,
	}
}

func (s *TCP) Server() {
	var err error
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	log.Println("TCP server starting...")
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("Error while listening tcp connection: %v\n", err)
	}
	defer s.listener.Close()
	log.Println("TCP listener started.")

	s.conn, err = s.listener.Accept()
	s.io = s.conn
	if err != nil {
		log.Panicf("Error while accepting tcp connection: %v\n", err)
	}
	defer s.conn.Close()
	log.Println("TCP connection accepted.")
	log.Println("TCP server started.")
	close(s.chans.ConnectChan)

	go s.sendLoop()
	go s.receiveLoop()

	<-s.chans.QuitChan
	log.Println("TCP server closed.")
}

func (s *TCP) Client() {
	var err error
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	log.Println("TCP client starting...")
	s.conn, err = net.Dial("tcp", addr)
	s.io = s.conn
	if err != nil {
		log.Panicf("Error while dialing tcp connection: %v\n", err)
	}
	defer s.conn.Close()
	log.Println("TCP connection dialed.")
	log.Println("TCP client started.")
	close(s.chans.ConnectChan)

	go s.sendLoop()
	go s.receiveLoop()

	<-s.chans.QuitChan
	log.Println("TCP client closed.")
}
