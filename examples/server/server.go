package main

import (
	"fmt"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
)

//ServerInfo is a struct to server info
type ServerInfo struct {
	//Host
	Host string
	//PortNumber of the server
	PortNumber uint16
	//MaxWindow is the max sliding-window size
	MaxWindow string
}

var (
	p, seedValue float32
)

func recieve(servConn *net.UDPConn) {

}

func main() {
	servAddr, err := net.ResolveUDPAddr("udp", ":10001")
	fmt.Println(servAddr)
	errors.CheckError(err)

	servConn, err := net.ListenUDP("udp", servAddr)
	errors.CheckError(err)

	defer servConn.Close()

	// go read from the connection
	recieve(servConn)
}
