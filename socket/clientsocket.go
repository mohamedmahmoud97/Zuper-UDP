package socket

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	"github.com/vmihailenco/msgpack"
)

var (
	//RecData is the buffer of all data received
	RecData           = make([]byte, pckNo*512)
	lastAck     int32 = -1
	buffer            = make(map[uint32][]byte)
	corruptProb float32
	fileName    string
	pckNo       uint16
)

//CreateClientSocket in client-side
func CreateClientSocket(localAddr, servAddr *net.UDPAddr) *net.UDPConn {
	conn, err := net.DialUDP("udp", localAddr, servAddr)
	errors.CheckError(err)
	return conn
}

//SendToServer the filename of the needed file
func SendToServer(conn *net.UDPConn, window int, filename string) {
	fmt.Printf("hello the client is requesting file %v from server ... \n", filename)
	fileName = filename

	file := []byte(filename)

	noOfBytes := uint16(len(file))
	reqPacket := Packet{Data: file, pckNo: 1, Len: noOfBytes}

	b, err := msgpack.Marshal(&reqPacket)
	if err != nil {
		panic(err)
	}

	fmt.Println("Encoded the message ...")

	_, err = conn.Write(b)
	errors.CheckError(err)
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, packet *Packet) {
	ack := AckPacket{Seqno: packet.Seqno}

	b, err := msgpack.Marshal(&ack)
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(b)
	errors.CheckError(err)
}

//ReceiveFromServer any ack packet
func ReceiveFromServer(conn *net.UDPConn, algo string) {

	//Read data from server socket
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

		if algo == "sw" {
			go sendResponse(conn, addr, &packet)
			appendFile(packet.Data)
			checkOnPck(&packet)
		} else if algo == "gbn" {
			if int32(packet.Seqno) == lastAck+1 {
				lastAck = int32(packet.Seqno)
				fmt.Printf("last ack packet is %v\n", lastAck)
				appendFile(packet.Data)
				go sendResponse(conn, addr, &packet)
				checkOnPck(&packet)
			} else if int32(packet.Seqno) > lastAck+1 && lastAck != -1 {
				//change seqno of ack packet to last delivered packet
				packet.Seqno = uint32(lastAck)
				go sendResponse(conn, addr, &packet)
			} else if int32(packet.Seqno) > lastAck+1 && lastAck == -1 {

			}
		} else if algo == "sr" {
			buffer[packet.Seqno] = packet.Data
			appendFile(packet.Data)
			go sendResponse(conn, addr, &packet)
			checkOnPck(&packet)
		}

	}
}

///append to buffer to build file later on
func appendFile(data []byte) {
	for i := 0; i < len(data); i++ {
		RecData = append(RecData, data[i])
	}
}

//build the requested file at the client-side
func buildFile() {
	fmt.Println("Building File ... ")
	err := ioutil.WriteFile(fileName, RecData, 0644)
	errors.CheckError(err)
	fmt.Println("Finished ... ")
	os.Exit(0)
}

func checkOnPck(packet *Packet) {
	if packet.Seqno == 0 {
		pckNo = packet.pckNo
	} else if int(packet.Seqno) == int(pckNo)-1 {
		buildFile()
	}
}
