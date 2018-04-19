package socket

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	"github.com/vmihailenco/msgpack"
)

var (
	//AckCheck is a channel for receiving seqno of ack packets
	AckCheck = make(chan uint32)
)

//CreateSerSocket in the server-side
func CreateSerSocket(servAddr *net.UDPAddr) *net.UDPConn {
	//create the socket on the port number
	servConn, err := net.ListenUDP("udp", servAddr)
	errors.CheckError(err)
	return servConn
}

func encodeFile(fileName string) []byte {
	var file bytes.Buffer
	file.WriteString("/home/mohamedmahmoud/Workspaces/Zuper-UDP/")
	file.WriteString(fileName)
	dat, err := ioutil.ReadFile(file.String())
	errors.CheckError(err)
	return dat
}

func sendToClient(conn *net.UDPConn, window int, addr *net.UDPAddr, algo, filename string, plp float32) {
	//read file into bytes
	dataBytes := encodeFile(filename)

	var seqNum uint32
	size := 512
	previous := 0
	r := 512

	noChunk := float64(len(dataBytes)) / float64(r)
	noChunks := uint32(math.Ceil(noChunk))

	packets := []Packet{}

	//make data packets and segment the file to be sent and assign seqNumber
	for seqNum < noChunks {
		chunk := dataBytes[previous:r]
		noOfBytes := uint16(len(chunk))
		piko := Packet{}
		piko.Data = chunk
		piko.Len = noOfBytes
		piko.Seqno = seqNum
		piko.Cksum = uint16(noChunks)
		packets = append(packets, piko)
		//making packets
		fmt.Printf("making packet %d ...\n", seqNum)
		seqNum++
		previous += size
		r += size
	}

	noOfChunks := int(noChunks)

	//send the packets in the way of the given algo
	reliableSend(packets, noOfChunks, conn, window, addr, algo, plp)
}

func reliableSend(packets []Packet, noChunks int, conn *net.UDPConn, window int, addr *net.UDPAddr, algo string, plp float32) {
	switch algo {
	case "sw":
		SW(packets, noChunks, conn, addr, plp)
	case "gbn":
		GBN(packets, noChunks, conn, addr, window, plp)
	case "sr":
		SR(packets, noChunks, conn, addr, window, plp)
	}
	fmt.Print("finished ... \n")
}

func sendAckToClient(conn *net.UDPConn, addr *net.UDPAddr) {
	ack := AckPacket{Seqno: 0}

	b, err := msgpack.Marshal(&ack)
	if err != nil {
		panic(err)
	}

	_, err = conn.WriteToUDP(b, addr)
	errors.CheckError(err)
}

//ReceiveReqFromClients any packet
func ReceiveReqFromClients(conn *net.UDPConn, buf []byte, length int, addr *net.UDPAddr, windowSize int, algo string, plp float32) {
	var packet Packet

	err := msgpack.Unmarshal(buf, &packet)
	if err != nil {
		panic(err)
	}

	n := len(packet.Data)
	filename := string(packet.Data[:n])
	fmt.Printf("requested the filename: %v", filename)

	sendAckToClient(conn, addr)
	sendToClient(conn, windowSize, addr, algo, filename, plp)
}

//ReceiveAckFromClients any packet
func ReceiveAckFromClients(conn *net.UDPConn, buf []byte, length int, addr *net.UDPAddr, windowSize int, algo string) {
	var packet AckPacket

	err := msgpack.Unmarshal(buf, &packet)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Received Ack of packet with seqno %v \n", packet.Seqno)

	//a channel for sending seqno
	AckCheck <- packet.Seqno
}
