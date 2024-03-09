package main

import (
	"flag"

	"github.com/fevzibabaoglu/InternetProtocolComparison/protocols"
)

func main() {
	var ip string
	var protocol string
	var isServer bool

	flag.StringVar(&ip, "ip", "localhost", "IP")
	flag.StringVar(&protocol, "protocol", "tcp", "Protocol")
	flag.BoolVar(&isServer, "server", false, "Enable server mode")
	flag.Parse()

	switch isServer {
	case true:
		switch protocol {
		case "tcp":
			protocols.TCPServer(ip)
		case "udp":
			protocols.UDPServer(ip)
		case "quic":
			protocols.QUICServer(ip)
		}
	case false:
		switch protocol {
		case "tcp":
			protocols.TCPClient(ip)
		case "udp":
			protocols.UDPClient(ip)
		case "quic":
			protocols.QUICClient(ip)
		}
	}
}
