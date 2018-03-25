package main

import (
	"fmt"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	socket "github.com/mohamedmahmoud97/Zuper-UDP/socket"
)

//ServerInfo is a struct to server info
type ServerInfo struct {
	//PortNumber of the server
	PortNumber string
	//MaxWindow is the max sliding-window size
	MaxWindow int
}

var (
	p, seedValue float32
)

func main() {
	server := ServerInfo{":10001", 5}
	p = 0.1
	seedValue = 50
	servAddr, err := net.ResolveUDPAddr("udp", server.PortNumber)
	fmt.Printf("connection in server on port %v", servAddr)
	errors.CheckError(err)

	//create the socket in server-side
	servConn := socket.CreateSerSocket(servAddr)

	// go read from the connection
	for {
		socket.ReceiveFromClients(servConn, p, seedValue)
	}
}
