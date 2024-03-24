package comms

import (
	"crypto/rand"
	"log"
	"strings"
	"sync"
)

const (
	messageSize  = 128 // Max 'messageSize' is '255' due to usage of 1 byte uint8 as 'byteSize' in '(*socket).Receive()'.
	messageCount = 1000000
)

func Run(protocol, ip, port string, bidirectional bool) {
	var s1, s2 ISocket
	var wg sync.WaitGroup

	// create sockets
	switch strings.ToUpper(protocol) {
	case "TCP":
		s1 = NewTCP(ip, port)
		s2 = NewTCP(ip, port)
	case "UDP":
		s1 = NewUDP(ip, port)
		s2 = NewUDP(ip, port)
	case "QUIC":
		s1 = NewQUIC(ip, port)
		s2 = NewQUIC(ip, port)
	}

	// start sockets
	log.Printf("%s test starting.\n", protocol)
	wg.Add(1)
	go func() {
		defer wg.Done()
		s1.Server()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s2.Client()
	}()
	wg.Wait()

	// defer close sockets
	defer s1.Close()
	defer s2.Close()

	// generating random strings
	var messages [][]byte
	for i := 0; i < messageCount; i++ {
		message := make([]byte, messageSize)
		rand.Read(message)
		messages = append(messages, message)
	}

	// server-send/client-receive test
	if bidirectional {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, message := range messages {
				s1.Send(string(message))
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, message := range messages {
				messageReceived := s2.Receive()
				if string(message) != messageReceived {
					log.Fatalln("client receive failed.")
				}
			}
		}()
	}

	// client-send/server-receive test
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, message := range messages {
			s2.Send(string(message))
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, message := range messages {
			messageReceived := s1.Receive()
			if string(message) != messageReceived {
				log.Fatalln("server receive failed.")
			}
		}
	}()

	// done
	wg.Wait()
	log.Printf("%s successful.\n\n", protocol)
}
