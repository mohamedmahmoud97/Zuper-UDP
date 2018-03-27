package socket

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
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

func reliableSend(packets []Packet, noChunks int, conn *net.UDPConn, window int) {
	// for i := 0; i < noChunks; i++ {
	// 	fmt.Fprintf(conn, packets[i])

	// 	_, err := bufio.NewReader(conn).Read(p)
	// 	if err == nil {
	// 		fmt.Printf("%s\n", p)
	// 	} else {
	// 		fmt.Printf("Some error %v\n", err)
	// 	}
	// }

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	for i := 0; i < noChunks; i++ {
		encoder.Encode(packets[i])
		conn.Write(buffer.Bytes())
		buffer.Reset()
	}
}

//SendToServer packets
func SendToServer(conn *net.UDPConn, window int, filename string) {
	fmt.Println("hello the client is starting the sending process ...")

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
		piko := Packet{data: chunk, len: noOfBytes, seqno: seqNum}
		packets = append(packets, piko)
		//making packets
		fmt.Printf("making packets: \nPacket %d", seqNum)
		fmt.Print(packets[seqNum])
		seqNum++
		previous += size
		r += size
	}

	noOfChunks := int(noChunks)
	reliableSend(packets, noOfChunks, conn, window)
}
