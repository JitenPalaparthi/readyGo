package messaging

import (
	"github.com/golang/glog"
	nats "github.com/nats-io/nats.go"
)

// Messaging is to maintain messagig information
type Messaging struct {
	ChanMessage chan Message
	ChanError   chan error
	Client      interface{}
}

// Message is a message content, ideally for a channel
type Message struct {
	Data    []byte
	Subject string
}

// Disconnect is to disconnect the given connection
func (m *Messaging) Disconnect() {
	if m.Client != nil {
		m.Client.(*nats.Conn).Close()
	}
}

// Init is to Initiate the messaging stuff
func (m *Messaging) Init() {
	if m.ChanMessage == nil {
		m.ChanMessage = make(chan Message, 20)
		go m.Publish()
	}
}

// Publish is to  produce messages by the producer
func (m *Messaging) Publish() {
	for msg := range m.ChanMessage {
		if m.Client != nil {
			m.Client.(*nats.Conn).Publish(msg.Subject, msg.Data)
		}
	}
}

// Subscribe is to  subscribe messages by the producer
func (m *Messaging) Subscribe() {
	for msg := range m.ChanMessage {
		if m.Client != nil {
			// Write subsctibe logic here
			glog.Info(msg) // Just to make sure msg does not give any error
			//m.Client.(*nats.Connection).Publish(msg.Subject, msg.Data)
		}
	}
}
