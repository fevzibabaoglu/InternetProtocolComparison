package main

import (
	"flag"

	"github.com/fevzibabaoglu/InternetProtocolComparison/comms"
)

func main() {
	var ip string
	var port string

	flag.StringVar(&ip, "ip", "localhost", "IP")
	flag.StringVar(&port, "port", "3000", "Port")
	flag.Parse()

	comms.Run("TCP", ip, port, true)
	comms.Run("UDP", ip, port, false)
	comms.Run("QUIC", ip, port, true)
}
