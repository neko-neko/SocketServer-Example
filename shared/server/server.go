package server

import (
	"net"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"sync"

	"github.com/neko-neko/SocketServer-Example/shared/log"
)

// connection id
var connectionId uint64

// Server
type Server struct {
	// TCP Address
	addr *net.TCPAddr

	// Socket
	socket *net.TCPListener

	// Connections
	connections map[*Connection]struct{}

	// Receive channel
	packetQueue chan PacketQueue

	// Signal
	signal chan os.Signal

	// Mutex
	mutex *sync.Mutex
}

// packet queue
type PacketQueue struct {
	Packet     []byte
	Connection *Connection
}

// Create new socket server
func NewServer(host string, port int, packetQueue chan PacketQueue) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	log.Info("Starting server...", addr.String())

	return &Server{
		packetQueue: packetQueue,
		addr:        addr,
		socket:      nil,
		connections: make(map[*Connection]struct{}),
		signal:      make(chan os.Signal, 1),
		mutex:       new(sync.Mutex),
	}, nil
}

// Listen TCP connection
func (s *Server) listen() error {
	l, err := net.ListenTCP("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Info("Listen server...")

	s.socket = l
	return nil
}

// Accept connection
func (s *Server) handleClient() {
	for {
		conn, err := s.socket.Accept()
		if err != nil {
			log.Warn("Accept error.", err)
			continue
		}

		log.Debug("Accepted connection.", conn.RemoteAddr())

		c := NewConnection(conn)
		s.connections[c] = struct{}{}
		c.TransitionState(Connected)
		c.Id = s.generateId()

		// Connection tasks
		c.OnClosed = s.closeConnectionFunc(c)
		c.Wait(s.packetQueue)
	}
}

// signal handler
func (s *Server) handleSignal() {
	signal.Notify(s.signal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	for sig := range s.signal {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL:
			// close all connections
			s.close()

			// close socket
			s.socket.Close()

			os.Exit(0)
		}
	}
}

// close all connections
func (s *Server) close() {
	for c := range s.connections {
		c.Closing <- struct{}{}
	}
}

// close connection callback
func (s *Server) closeConnectionFunc(connection *Connection) func() {
	return func() {
		delete(s.connections, connection)
	}
}

// run server
func (s *Server) Run() error {
	// listen port
	err := s.listen()
	if err != nil {
		return err
	}

	// accept connections and manage stream
	go s.handleClient()

	// signal handler
	go s.handleSignal()

	log.Info("Server is running!")

	return nil
}

// Notify message
func (s *Server) Notify(conn *Connection, message []byte) {
	if _, ok := s.connections[conn]; ok {
		conn.Message <- message
	}
}

// broadcast message
func (s *Server) NotifyAll(message []byte) {
	for c := range s.connections {
		c.Message <- message
	}
}

// generate connection id
func (s *Server) generateId() uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	connectionId++

	return connectionId
}
