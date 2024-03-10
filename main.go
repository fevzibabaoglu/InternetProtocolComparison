package main

import (
	"flag"

	"github.com/fevzibabaoglu/InternetProtocolComparison/comms"
)

func main() {
	var protocol string
	var ip string
	var port string
	var isServer bool

	flag.StringVar(&protocol, "protocol", "tcp", "Select protocol [TCP/UDP/QUIC]")
	flag.StringVar(&ip, "ip", "localhost", "IP")
	flag.StringVar(&port, "port", "3000", "Port")
	flag.BoolVar(&isServer, "server", false, "Enable server mode")
	flag.Parse()

	comms.Run(protocol, ip, port, isServer)
}
