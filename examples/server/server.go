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

	//create the socket on the port number
	servConn, err := net.ListenUDP("udp", servAddr)
	errors.CheckError(err)

	//defer servConn.Close()

	// go read from the connection
	for{
		socket.ReceiveFromClients(servConn, p, seedValue)
	}
}
