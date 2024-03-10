package comms

const (
	channelSize = 10
)

type Channels struct {
	SendChan    chan []byte
	ReceiveChan chan []byte
	ConnectChan chan struct{}
	QuitChan    chan struct{}
}

func newChannels() Channels {
	return Channels{
		SendChan:    make(chan []byte, channelSize),
		ReceiveChan: make(chan []byte, channelSize),
		ConnectChan: make(chan struct{}),
		QuitChan:    make(chan struct{}),
	}
}
