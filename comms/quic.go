package comms

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"

	"github.com/quic-go/quic-go"
)

type QUIC struct {
	socket
	protos   []string
	listener *quic.Listener
	conn     quic.Connection
	stream   quic.Stream
}

func newQUIC(ip string, port string) *QUIC {
	return &QUIC{
		socket: socket{
			ip:    ip,
			port:  port,
			io:    nil,
			chans: newChannels(),
		},
		protos:   []string{"quic"},
		listener: nil,
		conn:     nil,
		stream:   nil,
	}
}

func (s *QUIC) Server() {
	var err error
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	log.Println("QUIC server starting...")
	s.listener, err = quic.ListenAddr(addr, s.serverTLSConfig(), nil)
	if err != nil {
		log.Panicf("Error while listening quic connection: %v\n", err)
	}
	defer s.listener.Close()
	log.Println("QUIC listener started.")

	s.conn, err = s.listener.Accept(context.Background())
	if err != nil {
		log.Panicf("Error while accepting quic connection: %v\n", err)
	}
	log.Println("QUIC connection accepted.")

	s.stream, err = s.conn.AcceptStream(context.Background())
	s.io = s.stream
	if err != nil {
		log.Panicf("Error while accepting quic stream: %v\n", err)
	}
	defer s.stream.Close()
	s.receiveBlankInitial()
	log.Println("QUIC stream accepted.")
	log.Println("QUIC server started.")
	close(s.chans.ConnectChan)

	go s.sendLoop()
	go s.receiveLoop()

	<-s.chans.QuitChan
	log.Println("QUIC server closed.")
}

func (s *QUIC) Client() {
	var err error
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	log.Println("QUIC client starting...")
	s.conn, err = quic.DialAddr(context.Background(), addr, s.clientTLSConfig(), nil)
	if err != nil {
		log.Panicf("Error while dialing quic connection: %v\n", err)
	}
	log.Println("QUIC connection dialed.")

	s.stream, err = s.conn.OpenStreamSync(context.Background())
	s.io = s.stream
	if err != nil {
		log.Panicf("Error while opening quic stream: %v\n", err)
	}
	defer s.stream.Close()
	s.sendBlankInitial()
	log.Println("QUIC stream opened.")
	log.Println("QUIC client started.")
	close(s.chans.ConnectChan)

	go s.sendLoop()
	go s.receiveLoop()

	<-s.chans.QuitChan
	log.Println("QUIC client closed.")
}

func (socket *QUIC) receiveBlankInitial() {
	buffer := make([]byte, 1)
	_, _ = socket.stream.Read(buffer)
}

func (socket *QUIC) sendBlankInitial() {
	for {
		_, err := socket.stream.Write(make([]byte, 1))
		if err == nil {
			break
		}
	}
}

func (socket *QUIC) serverTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   socket.protos,
	}
}

func (socket *QUIC) clientTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         socket.protos,
	}
}
