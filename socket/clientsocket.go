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

//CreateClientSocket in client-side
func CreateClientSocket(localAddr, servAddr *net.UDPAddr) *net.UDPConn {
	conn, err := net.DialUDP("udp", localAddr, servAddr)
	errors.CheckError(err)
	return conn
}

func encodeFile(fileName string) []byte {
	var file bytes.Buffer
	file.WriteString("/home/mohamedmahmoud/Workspaces/Zuper-UDP/")
	file.WriteString(fileName)
	dat, err := ioutil.ReadFile(file.String())
	errors.CheckError(err)
	return dat
}

func reliableSend(packets []Packet, noChunks int, conn *net.UDPConn, window int, servAddr *net.UDPAddr) {
	for i := 0; i < noChunks; i++ {
		b, err := msgpack.Marshal(&packets[i])
		if err != nil {
			panic(err)
		}
		_, err = conn.Write(b)
	}
}

//SendToServer packets
func SendToServer(conn *net.UDPConn, window int, filename string, servAddr *net.UDPAddr) {
	fmt.Print("hello the client is starting the sending process ... \n")

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
		piko := Packet{Data: chunk, Len: noOfBytes, Seqno: seqNum}
		packets = append(packets, piko)
		//making packets
		fmt.Printf("making packet %d ...\n", seqNum)
		seqNum++
		previous += size
		r += size
	}

	noOfChunks := int(noChunks)
	reliableSend(packets, noOfChunks, conn, window, servAddr)
}

//ReceiveFromServer any ack packet
func ReceiveFromServer(conn *net.UDPConn) {
	buf := make([]byte, 600)
	length, _, err := conn.ReadFromUDP(buf[0:])
	errors.CheckError(err)

	if length < 30 && length != 0 {
		var packet Packet

		err := msgpack.Unmarshal(buf, &packet)
		if err != nil {
			panic(err)
		}

		fmt.Println(packet.Cksum)
	}

}
