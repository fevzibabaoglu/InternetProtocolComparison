package comms

import (
	"fmt"
	"log"
	"net"
)

type UDP struct {
	socket
	udpAddr *net.UDPAddr
	conn    *net.UDPConn
}

func NewUDP(ip string, port string) *UDP {
	return &UDP{
		socket: socket{
			ip:   ip,
			port: port,
			io:   nil,
		},
		udpAddr: nil,
		conn:    nil,
	}
}

func (s *UDP) Close() {
	s.conn.Close()
}

func (s *UDP) Server() {
	var err error
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	s.udpAddr, err = net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Panicf("Error while resolving udp address: %v\n", err)
	}

	s.conn, err = net.ListenUDP("udp", s.udpAddr)
	s.io = s.conn
	if err != nil {
		log.Panicf("Error while listening udp connection: %v\n", err)
	}
}

func (s *UDP) Client() {
	var err error
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	s.udpAddr, err = net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Panicf("Error while resolving udp address: %v\n", err)
	}

	s.conn, err = net.DialUDP("udp", nil, s.udpAddr)
	s.io = s.conn
	if err != nil {
		log.Panicf("Error while dialing udp connection: %v\n", err)
	}
}
