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

func newUDP(ip string, port string) *UDP {
	return &UDP{
		socket: socket{
			ip:    ip,
			port:  port,
			io:    nil,
			chans: newChannels(),
		},
		udpAddr: nil,
		conn:    nil,
	}
}

func (s *UDP) Server() {
	var err error
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	log.Println("UDP server starting...")
	s.udpAddr, err = net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Panicf("Error while resolving udp address: %v\n", err)
	}
	log.Println("UDP address resolved.")

	s.conn, err = net.ListenUDP("udp", s.udpAddr)
	s.io = s.conn
	if err != nil {
		log.Panicf("Error while listening udp connection: %v\n", err)
	}
	defer s.conn.Close()
	log.Println("UDP listener started.")
	log.Println("UDP server started.")
	close(s.chans.ConnectChan)

	// go s.sendLoop()
	go s.receiveLoop()

	<-s.chans.QuitChan
	log.Println("UDP server closed.")
}

func (s *UDP) Client() {
	var err error
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	log.Println("UDP client starting...")
	s.udpAddr, err = net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Panicf("Error while resolving udp address: %v\n", err)
	}
	log.Println("UDP address resolved.")

	s.conn, err = net.DialUDP("udp", nil, s.udpAddr)
	s.io = s.conn
	if err != nil {
		log.Panicf("Error while dialing udp connection: %v\n", err)
	}
	defer s.conn.Close()
	log.Println("UDP connection dialed.")
	log.Println("UDP client started.")
	close(s.chans.ConnectChan)

	go s.sendLoop()
	// go s.receiveLoop()

	<-s.chans.QuitChan
	log.Println("UDP client closed.")
}
