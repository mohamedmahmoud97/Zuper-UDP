package socket

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	"github.com/mohamedmahmoud97/Zuper-UDP/socket"
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

func reliableSend(packets []sockets.Packet){
	for i := 0; i < noChunks; i++ {
		fmt.Fprintf(conn, packets[i])

		_, err := bufio.NewReader(conn).Read(p)
		if err == nil {
			fmt.Printf("%s\n", p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
	}
}

//SendToServer packets
func SendToServer(conn *net.UDPConn, window int, filename string) {
	p := make([]byte, 2048)

	//read file into bytes
	dataBytes := encodeFile(filename)

	var previous, seqNum = 0
	r := 512

	noChunks := float64(dataBytes) / float64(r)
	int(math.Ceil(noChunks))

	packets := []socket.Packet

	//make data packets and segment the file to be sent and assign seqNumber
	for seqNum<noChunks { 
		chunk := dataBytes[previous:r]
		piko := socket.Packet{data: chunk, len: len(chunk), seqno: seqNum}
		packets.append(piko)
		seqNum += 1
		previous += r
	}

	reliableSend(packets)
}
