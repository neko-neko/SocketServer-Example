package main

import (
	"os"
	"strconv"
	"fmt"
	"time"

	"github.com/neko-neko/SocketServer-Example/shared/log"
	"github.com/neko-neko/SocketServer-Example/server"
)

const MaxPacketQueueSize = 1024

func main() {
	var worldUpdateCount int
	host := os.Getenv("SOCKET_SERVER_HOST")
	port, _ := strconv.Atoi(os.Getenv("SOCKET_SERVER_PORT"))

	log.SetLevel()
	log.SetOutput()

	// prepare server
	ch := make(chan server.PacketQueue, MaxPacketQueueSize)
	s, err := server.NewServer(host, port, ch)
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	// boot server
	rerr := s.Run()
	if rerr != nil {
		log.Panic(rerr)
		os.Exit(1)
	}

	// main loop
	timer := time.Tick(100 * time.Millisecond)
	for _ = range timer {
		select {
		// Receive Packet
		case p := <-ch:
			// response message
			s.Notify(p.Connection, []byte("your message received\n"))

			// broadcast message
			s.NotifyAll(p.Packet)
		// Update world
		default:
			worldUpdateCount++
			s.NotifyAll([]byte([]byte(fmt.Sprintf("worldUpdateCount is now %d", worldUpdateCount))))
		}
	}
}
