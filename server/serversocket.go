package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	socket "github.com/mohamedmahmoud97/Zuper-UDP/socket"
	"github.com/vmihailenco/msgpack"
)

var (
	//AckCheck is a channel for receiving seqno of ack packets
	AckCheck = make(chan uint32)
	flogS    *os.File
)

//CreateSerSocket in the server-side
func CreateSerSocket(servAddr *net.UDPAddr, flog *os.File) *net.UDPConn {
	flogS = flog

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

	packets := []socket.Packet{}

	//make data packets and segment the file to be sent and assign seqNumber
	for seqNum < noChunks {
		chunk := dataBytes[previous:r]
		noOfBytes := uint16(len(chunk))
		piko := socket.Packet{}
		piko.Data = chunk
		piko.Len = noOfBytes
		piko.Seqno = seqNum
		piko.PckNo = uint16(noChunks)
		packets = append(packets, piko)
		seqNum++
		previous += size
		r += size
	}
	log.SetOutput(flogS)
	log.Printf("Made %v packets ... \n", noChunks)

	noOfChunks := int(noChunks)

	//send the packets in the way of the given algo
	reliableSend(packets, noOfChunks, conn, window, addr, algo, plp)
}

func reliableSend(packets []socket.Packet, noChunks int, conn *net.UDPConn, window int, addr *net.UDPAddr, algo string, plp float32) {
	switch algo {
	case "sw":
		SW(packets, noChunks, conn, addr, plp)
	case "gbn":
		GBN(packets, noChunks, conn, addr, window, plp)
	case "sr":
		SR(packets, noChunks, conn, addr, window, plp)
	}
	log.SetOutput(flogS)
	log.Println("Finished ... ")
	fmt.Print("finished ... \n")
}

func sendAckToClient(conn *net.UDPConn, addr *net.UDPAddr) {
	ack := socket.AckPacket{Seqno: 0}

	b, err := msgpack.Marshal(&ack)
	if err != nil {
		panic(err)
	}

	_, err = conn.WriteToUDP(b, addr)
	errors.CheckError(err)
}

//ReceiveReqFromClients any packet
func ReceiveReqFromClients(conn *net.UDPConn, buf []byte, length int, addr *net.UDPAddr, windowSize int, algo string, plp float32) {
	var packet socket.Packet

	err := msgpack.Unmarshal(buf, &packet)
	if err != nil {
		panic(err)
	}

	n := len(packet.Data)
	filename := string(packet.Data[:n])
	log.SetOutput(flogS)
	log.Printf("A client requested filename: %v \n", filename)
	fmt.Printf("requested the filename: %v", filename)

	sendAckToClient(conn, addr)
	sendToClient(conn, windowSize, addr, algo, filename, plp)
}

//ReceiveAckFromClients any packet
func ReceiveAckFromClients(conn *net.UDPConn, buf []byte, length int, addr *net.UDPAddr, windowSize int, algo string) {
	var packet socket.AckPacket

	err := msgpack.Unmarshal(buf, &packet)
	if err != nil {
		panic(err)
	}

	log.SetOutput(flogS)
	log.Printf("Received Ack of packet with seqno %v \n", packet.Seqno)
	fmt.Printf("Received Ack of packet with seqno %v \n", packet.Seqno)

	//a channel for sending seqno
	AckCheck <- packet.Seqno
}
