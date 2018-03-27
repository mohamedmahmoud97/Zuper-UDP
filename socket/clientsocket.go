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
	// for i := 0; i < 1; i++ {
	// 	fmt.Fprintf(conn, "hey it's a fuckin test yo bitch be cool !!")

	// 	// _, err := bufio.NewReader(conn).Read(p)
	// 	// if err == nil {
	// 	// 	fmt.Printf("%s\n", p)
	// 	// } else {
	// 	// 	fmt.Printf("Some error %v\n", err)
	// 	// }
	// }

	// var buffer bytes.Buffer
	// encoder := gob.NewEncoder(&buffer)
	for i := 0; i < noChunks; i++ {
		// jsonRequest, err := json.Marshal(packets[i])
		// if err != nil {
		// 	log.Print("Marshal connection information failed.")
		// 	log.Fatal(err)
		// }

		// _, err := conn.Write([]byte(fmt.Sprintf("%v", packets[i])))
		// errors.CheckError(err)

		// var buf bytes.Buffer
		// if err := gob.NewEncoder(&buf).Encode(packets[i]); err != nil {
		// 	errors.CheckError(err)
		// }
		// _, err := conn.WriteToUDP(buf.Bytes(), servAddr)
		// errors.CheckError(err)

		// fmt.Printf("\n\n\nsending packet %d\n", i)
		// encoder.Encode(packets[25])
		// fmt.Println(buffer)
		// _, err := conn.Write(buffer.Bytes())
		// errors.CheckError(err)
		// buffer.Reset()

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
