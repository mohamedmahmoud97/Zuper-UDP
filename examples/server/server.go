package main

import (
	"fmt"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	socket "github.com/mohamedmahmoud97/Zuper-UDP/socket/serversocket"
)

//ServerInfo is a struct to server info
type ServerInfo struct {
	//PortNumber of the server
	PortNumber uint16
	//MaxWindow is the max sliding-window size
	MaxWindow string
}

var (
	p, seedValue float32
)

func main() {
	server := ServerInfo{":10001",5}
	p = 0.1
	seedValue = 
	servAddr, err := net.ResolveUDPAddr("up4", server.PortNumber)
	fmt.Printf("connection in server on port %v", servAddr)
	errors.CheckError(err)

	socket.CreateSocket(servAddr)

	// go read from the connection
	for{
		socket.ReceiveFromClients(servConn, p, seedValue)
	}
}
