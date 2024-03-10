package comms

import (
	"encoding/binary"
	"io"
	"log"
	"regexp"
	"strings"
)

type ISocket interface {
	Server()
	Client()
	GetChans() Channels
}

func NewISocket(protocol, ip, port string) ISocket {
	switch strings.ToUpper(protocol) {
	case "TCP":
		return newTCP(ip, port)
	case "UDP":
		return newUDP(ip, port)
	case "QUIC":
		return newQUIC(ip, port)
	default:
		return nil
	}
}

type socket struct {
	ip   string
	port string
	io   interface {
		io.Reader
		io.Writer
	}
	chans Channels
}

func (s *socket) GetChans() Channels {
	return s.chans
}

func (s *socket) sendLoop() {
	for {
		message := string(<-s.chans.SendChan)

		err := binary.Write(s.io, binary.LittleEndian, byte(len(message)))
		if err != nil {
			s.processStreamError(err)
		}

		_, err = s.io.Write([]byte(message))
		if err != nil {
			s.processStreamError(err)
		}
	}
}

func (s *socket) receiveLoop() {
	for {
		sizeByte := make([]byte, 1)
		_, err := s.io.Read(sizeByte)
		if err != nil {
			s.processStreamError(err)
		}

		buffer := make([]byte, sizeByte[0])
		_, err = s.io.Read(buffer)
		if err != nil {
			s.processStreamError(err)
		}

		s.chans.ReceiveChan <- buffer
	}
}

func (s *socket) processStreamError(errArg error) {
	// e.g. "INTERNAL_ERROR (local): write udp 127.0.0.1:3000->127.0.0.1:50291: use of closed network connection"

	pattern := regexp.MustCompile(`^.*: (.*)$`)
	matches := pattern.FindStringSubmatch(errArg.Error())
	regErrMsg := ""
	if matches != nil {
		regErrMsg = matches[1]
	}

	if errArg.Error() == "EOF" || regErrMsg == "use of closed network connection" {
		log.Fatalf("ERR: Connection closed\n")
	}

	log.Panicf("Error while interacting with stream (%s)\n", errArg)
}
