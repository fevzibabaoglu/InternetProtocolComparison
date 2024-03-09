package protocols

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func TCPServer(ip string) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Panicf("Error while listening tcp connection: %v\n", err)
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		log.Panicf("Error while accepting tcp connection: %v\n", err)
	}
	defer conn.Close()

	for i := 0; i < messageCount; i++ {
		sizeByte := make([]byte, 1)
		_, err := conn.Read(sizeByte)
		if err != nil {
			log.Panicf("Error while reading with tcp connection: %v\n", err)
		}

		buffer := make([]byte, sizeByte[0])
		_, err = conn.Read(buffer)
		if err != nil {
			log.Panicf("Error while reading with tcp connection: %v\n", err)
		}

		message := string(buffer)

		fmt.Printf("Received: %s\n", message)
	}
}

func TCPClient(ip string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Panicf("Error while dialing tcp connection: %v\n", err)
	}
	defer conn.Close()

	for i := 0; i < messageCount; i++ {
		message := fmt.Sprintf("Hello! %d", i)

		err := binary.Write(conn, binary.LittleEndian, byte(len(message)))
		if err != nil {
			log.Panicf("Error while adding endian to the message: %v\n", err)
			continue
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Panicf("Error while writing with tcp connection: %v\n", err)
			continue
		}

		fmt.Printf("Sending: %s\n", message)
	}
}
