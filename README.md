# SocketServer-Example
SocketServer-Example written in golang.  
Socket server has main-loop like MMO server.

## Installation
Multi-platform support
- Windows
- Mac OS X
- Linux

### Requirements
- [dep]("https://github.com/golang/dep")

### Build server and client

#### Download dependencies
```bash
$ dep ensure
```

#### Build Server
```bash
$ cd examples/server
$ go build main.go
```

#### Build Client
```bash
$ cd examples/client
$ go build main.go
```

## Run server and client
### Run server
```bash
$ SOCKET_SERVER_HOST=0.0.0.0 SOCKET_SERVER_PORT=11111 ./main
```

### Run client
```bash
$ SOCKET_SERVER_CONNECT_HOST=localhost SOCKET_SERVER_CONNECT_PORT=11111 ./main
```

## Credits
neko-neko

## License
MIT