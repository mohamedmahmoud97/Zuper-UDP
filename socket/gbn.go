package socket

import (
	"fmt"
	"net"

	"github.com/vmihailenco/msgpack"
)

// a map to give a value for all packets
// 0 = unsent, 1 = sent, 2 = sent and rcved ack
var (
	ackPack = make(map[int]int)
)

//GBN is the go-back-n algorithm
func GBN(packets []Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr, window int) {

	start := 0

	//make all the chunks have value 0 as unack
	for i := 0; i < noChunks; i++ {
		ackPack[i] = 0
	}

	//send the first packets with the window size
	sendWinPack(start, window, packets, conn, addr, noChunks)

	//loop until all the packets are sent and received their ack
	for (start + 3) < noChunks {
		ackpckt := int(<-AckCheck)

		if ackpckt == start {
			ackPack[ackpckt] = 2
			start = getNextStart(start, noChunks)
			if start != -1 {
				sendWinPack(start, window, packets, conn, addr, noChunks)
			}
		} else if ackPack[ackpckt] == 1 {
			ackPack[ackpckt] = 2
			nextUnSent := getNextStart(ackpckt, noChunks)
			if nextUnSent < start+4 && nextUnSent != -1 {
				sendWinPack(start, 1, packets, conn, addr, noChunks)
			}
		} else if ackPack[ackpckt] == 2 {
			sendWinPack(start, window, packets, conn, addr, noChunks)
		}
	}

	fmt.Print("finished ... ")

}

//send a package of packets with the size of the window
//used when start and after changing the start
func sendWinPack(start int, window int, packets []Packet, conn *net.UDPConn, addr *net.UDPAddr, noChunks int) {
	for i := start; i < start+window && i < noChunks; i++ {
		if ackPack[i] == 0 {
			b, err := msgpack.Marshal(&packets[i])
			if err != nil {
				panic(err)
			}
			_, err = conn.WriteToUDP(b, addr)
			ackPack[i] = 1
			fmt.Printf("Sent packet %v ... \n", i)
		}
	}
}

//get the next start which is not sent or not acked
func getNextStart(start int, noChunks int) int {
	for start < noChunks {
		start++
		if ackPack[start] == 0 || ackPack[start] == 1 {
			return start
		}
	}
	//return -1 if it finished sending all packets
	return -1
}
