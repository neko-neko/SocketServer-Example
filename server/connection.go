package server

import (
	"fmt"
	"io"
	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/neko-neko/SocketServer-Example/shared/log"
)

// Receive buffer size
const (
	MinBufferByteSize = 1
	MaxBufferByteSize = 1024
)

// Connection state
const (
	Initialize = iota
	Connected
	Closing
	Closed
)

// Unknown connection state
type UnknownConnectionStateError struct {
	State int
	Err   error
}

// Unknown connection state error
func (err *UnknownConnectionStateError) Error() string {
	return fmt.Sprintf("Unknown connection state %d.", err.State)
}

// Client connection socket
type Connection struct {
	// Connection ID
	Id uint64

	// send message channel
	Message chan []byte

	// Close channel
	Closing chan struct{}

	// Connection state
	state int

	// TCP connection
	conn net.Conn

	// Callback state transition
	OnConnected func()
	OnClosing   func()
	OnClosed    func()
}

// Create new connection
func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		Message: make(chan []byte),
		Closing: make(chan struct{}),
		state:   Initialize,
		conn:    conn,
	}
}

// wait events
func (c *Connection) Wait(packetQueue chan PacketQueue) {
	go c.recvWait(packetQueue)
	go c.sendWait()
	go c.closeWait()
}

// Transition connection state and run callback
func (c *Connection) TransitionState(state int) error {
	switch state {
	case Connected:
		if c.OnConnected != nil {
			c.OnConnected()
		}
		break
	case Closing:
		if c.OnClosing != nil {
			c.OnClosing()
		}
	case Closed:
		if c.OnClosed != nil {
			c.OnClosed()
		}
	default:
		err := &UnknownConnectionStateError{
			State: state,
		}
		return err
	}
	c.state = state

	return nil
}

// Read packet from socket
func (c *Connection) recvWait(packetQueue chan PacketQueue) {
	buf := make([]byte, MaxBufferByteSize)

	for {
		size, err := c.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				c.close()
				log.Debug("Disconnected connection.", c.conn.RemoteAddr())
			} else if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				log.Error(err)
			}
			break
		}

		if size < MinBufferByteSize {
			log.Warn("Not enough length of minimum packet size.", spew.Sprint(buf[:]))
			continue
		}

		if size > MaxBufferByteSize {
			log.Warn("Too long size packet.", spew.Sprint(buf[:]))
			continue
		}

		log.Debug(spew.Sprint(buf[:size]))

		packetQueue <- PacketQueue{
			Packet:     buf[:size],
			Connection: c,
		}
	}
}

// send packet to client
func (c *Connection) sendWait() {
	for {
		select {
		case m := <-c.Message:
			if c.state == Connected {
				c.conn.Write(m)
			}
			break
		default:
			// nothing to do
		}
	}
}

// close wait
func (c *Connection) closeWait() {
	<-c.Closing
	c.close()
}

// close connection
func (c *Connection) close() {
	c.TransitionState(Closing)
	c.conn.Close()
	c.TransitionState(Closed)
}
