package socket

import (
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
	ack := AckPacket{Cksum: 20}

	b, err := msgpack.Marshal(&ack)
	if err != nil {
		panic(err)
	}
	_, err = conn.WriteToUDP(b, addr)
	fmt.Println(ack.Cksum)
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
	sendResponse(conn, addr, &packet)
}
