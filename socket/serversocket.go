package socket

import (
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
)

var (
	//Data is the buffer of all data received
	Data = make([][512]byte, 512)
)

//CreateSerSocket in the server-side
func CreateSerSocket(servAddr *net.UDPAddr) *net.UDPConn {
	//create the socket on the port number
	servConn, err := net.ListenUDP("udp", servAddr)
	errors.CheckError(err)
	return servConn
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
	_, err := conn.WriteToUDP(buf, addr)
	errors.CheckError(err)
}

//ReceiveFromClients any packet
func ReceiveFromClients(conn *net.UDPConn) {
	var buf []byte

	_, addr, err := conn.ReadFromUDP(buf[0:])
	errors.CheckError(err)

	// Data.append(buf)
	if len(buf) > 0 {
		go sendResponse(conn, addr, buf)
	}
}
