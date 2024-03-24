package comms

import (
	"encoding/binary"
	"io"
	"log"
)

type ISocket interface {
	Server()
	Client()
	Close()
	Send(string)
	Receive() string
}

type socket struct {
	ip   string
	port string
	io   interface {
		io.Reader
		io.Writer
	}
}

func (s *socket) Send(message string) {
	err := binary.Write(s.io, binary.LittleEndian, byte(len(message)))
	if err != nil {
		log.Fatalf("ERROR WHILE WRITING: %v\n", err)
	}

	_, err = s.io.Write([]byte(message))
	if err != nil {
		log.Fatalf("ERROR WHILE WRITING: %v\n", err)
	}
}

func (s *socket) Receive() string {
	sizeByte := make([]byte, 1)
	_, err := s.io.Read(sizeByte)
	if err != nil {
		log.Fatalf("ERROR WHILE READING: %v\n", err)
	}

	buffer := make([]byte, sizeByte[0])
	totalRead := 0

	for totalRead != int(sizeByte[0]) {
		n, err := s.io.Read(buffer[totalRead:])
		totalRead += n
		if err != nil {
			log.Fatalf("ERROR WHILE READING: %v\n", err)
		}
	}

	return string(buffer)
}
