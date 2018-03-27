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
	fmt.Println(packet.Data)
	errors.CheckError(err)
}

//ReceiveFromClients any packet
func ReceiveFromClients(conn *net.UDPConn, buf []byte, length int, addr *net.UDPAddr) {
	var packet Packet

	// fmt.Println("\nStarting to accept from clients ... ")

	// if err := gob.NewDecoder(bytes.NewReader(buf[:length])).Decode(&packet); err != nil {
	// 	errors.CheckError(err)
	// }

	err := msgpack.Unmarshal(buf, &packet)
	if err != nil {
		panic(err)
	}

	// err := json.Unmarshal(buf[:length], &packet)
	// if err != nil {
	// 	fmt.Print("Unmarshal server response failed.")
	// 	log.Fatal(err)
	// }

	// buffer := bytes.NewBuffer(buf[:length])
	// decoder := gob.NewDecoder(buffer)
	// decoder.Decode(&packet)
	Data = append(Data, packet.Data)
	fmt.Println("am here yo old ass ....")
	fmt.Println(packet)
	//go sendResponse(conn, addr, &packet)
}
