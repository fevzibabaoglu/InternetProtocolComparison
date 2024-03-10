package comms

import (
	"fmt"
	"sync"
)

const (
	messageCount = 10
)

func Run(protocol, ip, port string, isServer bool) {
	iSocket := NewISocket(protocol, ip, port)

	switch isServer {
	case true:
		go iSocket.Server()
	case false:
		go iSocket.Client()
	}

	<-iSocket.GetChans().ConnectChan
	var wg sync.WaitGroup

	wg.Add(1)
	go func(iSocket ISocket) {
		defer wg.Done()

		for i := 0; i < messageCount; i++ {
			message := fmt.Sprintf("Hello! %d", i)
			iSocket.GetChans().SendChan <- []byte(message)
		}
	}(iSocket)

	wg.Add(1)
	go func(iSocket ISocket) {
		defer wg.Done()

		for i := 0; i < messageCount; i++ {
			message := <-iSocket.GetChans().ReceiveChan
			fmt.Println(string(message))
		}
	}(iSocket)

	wg.Wait()
	close(iSocket.GetChans().QuitChan)
}
