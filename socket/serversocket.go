package socket

import (
	"fmt"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
)

var (
	//Data is the buffer of all data received
	Data = make([][]byte, 1024)
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
	fmt.Println(string(buf))
	errors.CheckError(err)
}

//ReceiveFromClients any packet
func ReceiveFromClients(conn *net.UDPConn) {
	var buf []byte

	_, addr, err := conn.ReadFromUDP(buf[0:])
	errors.CheckError(err)

	if len(buf) > 0 {
		Data = append(Data, buf)
		fmt.Println(buf)
		go sendResponse(conn, addr, buf)
	}
}
