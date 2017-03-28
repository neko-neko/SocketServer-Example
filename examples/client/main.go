package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/neko-neko/SocketServer-Example/shared/log"
)

func main() {
	host := os.Getenv("SOCKET_SERVER_CONNECT_HOST")
	port, _ := strconv.Atoi(os.Getenv("SOCKET_SERVER_CONNECT_PORT"))

	log.SetLevel()
	log.SetOutput()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Panic(err)
		return
	}
	defer conn.Close()

	go handleMessage(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		conn.Write([]byte(scanner.Text()))
	}
}

func handleMessage(conn net.Conn) {
	for {
		readBuf := make([]byte, 1024)
		readLen, _ := conn.Read(readBuf)

		log.Info(string(readBuf[:readLen]))
	}
}
