package protocols

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"

	"github.com/quic-go/quic-go"
)

func QUICServer(ip string) {
	ln, err := quic.ListenAddr(fmt.Sprintf("%s:%d", ip, port), serverTLSConfig(), nil)
	if err != nil {
		log.Panicf("Error while listening quic connection: %v\n", err)
	}
	defer ln.Close()

	conn, err := ln.Accept(context.Background())
	if err != nil {
		log.Panicf("Error while accepting quic connection: %v\n", err)
	}

	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		log.Panicf("Error while accepting quic stream: %v\n", err)
	}
	defer stream.Close()

	for i := 0; i < messageCount; i++ {
		sizeByte := make([]byte, 1)
		_, err := stream.Read(sizeByte)
		if err != nil {
			log.Panicf("Error while reading with quic connection: %v\n", err)
		}

		buffer := make([]byte, sizeByte[0])
		_, err = stream.Read(buffer)
		if err != nil {
			log.Panicf("Error while reading with quic connection: %v\n", err)
		}

		message := string(buffer)

		fmt.Printf("Received: %s\n", message)
	}
}

func QUICClient(ip string) {
	conn, err := quic.DialAddr(context.Background(), fmt.Sprintf("%s:%d", ip, port), clientTLSConfig(), nil)
	if err != nil {
		log.Panicf("Error while dialing quic connection: %v\n", err)
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Panicf("Error while opening quic stream: %v\n", err)
	}
	defer stream.Close()

	for i := 0; i < messageCount; i++ {
		message := fmt.Sprintf("Hello! %d", i)

		err := binary.Write(stream, binary.LittleEndian, byte(len(message)))
		if err != nil {
			log.Panicf("Error while adding endian to the message: %v\n", err)
			continue
		}

		_, err = stream.Write([]byte(message))
		if err != nil {
			log.Panicf("Error while writing with tcp connection: %v\n", err)
			continue
		}

		fmt.Printf("Sending: %s\n", message)
	}
}

func serverTLSConfig() *tls.Config {
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
		NextProtos:   []string{"quic"},
	}
}

func clientTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic"},
	}
}
