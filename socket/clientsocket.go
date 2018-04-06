package socket

import (
	"fmt"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	"github.com/vmihailenco/msgpack"
)

var (
	//RecData is the buffer of all data received
	RecData = make([][]byte, 1024)
)

//CreateClientSocket in client-side
func CreateClientSocket(localAddr, servAddr *net.UDPAddr) *net.UDPConn {
	conn, err := net.DialUDP("udp", localAddr, servAddr)
	errors.CheckError(err)
	return conn
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, packet *Packet) {
	ack := AckPacket{Seqno: packet.Seqno, Cksum: 20}

	b, err := msgpack.Marshal(&ack)
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(b)
	errors.CheckError(err)
}

//SendToServer packets
func SendToServer(conn *net.UDPConn, window int, filename string) {
	fmt.Printf("hello the client is requesting file %v from server ... \n", filename)

	file := []byte(filename)

	noOfBytes := uint16(len(file))
	reqPacket := Packet{Data: file, Cksum: 20, Len: noOfBytes}

	b, err := msgpack.Marshal(&reqPacket)
	if err != nil {
		panic(err)
	}

	fmt.Println("Encoded the message ...")

	_, err = conn.Write(b)
	errors.CheckError(err)
}

//ReceiveFromServer any ack packet
func ReceiveFromServer(conn *net.UDPConn) {
	buf := make([]byte, 600)
	length, addr, err := conn.ReadFromUDP(buf[0:])
	errors.CheckError(err)

	if length > 0 {
		var packet Packet

		err := msgpack.Unmarshal(buf, &packet)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Delivered packet with seqno %v \n", packet.Seqno)

		RecData = append(RecData, packet.Data)
		go sendResponse(conn, addr, &packet)
	}

}
