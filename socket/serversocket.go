package socket

import (
	"bytes"
	"fmt"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	"github.com/vmihailenco/msgpack"
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

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, packet *Packet) {
	var ackStr bytes.Buffer
	ackStr.WriteString("delevired the packet with sequence number ")
	ackStr.WriteString(string(packet.Seqno))
	_, err := conn.WriteToUDP(ackStr.Bytes(), addr)
	//fmt.Println(packet.Data)
	fmt.Printf("delevired the packet with sequence number %v\n", packet.Seqno)
	errors.CheckError(err)
}

//ReceiveFromClients any packet
func ReceiveFromClients(conn *net.UDPConn, buf []byte, length int, addr *net.UDPAddr) {
	var packet Packet

	err := msgpack.Unmarshal(buf, &packet)
	if err != nil {
		panic(err)
	}

	Data = append(Data, packet.Data)
	go sendResponse(conn, addr, &packet)
}
