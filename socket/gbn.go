package socket

import (
	"fmt"
	"net"
)

// a map to give a value for all packets
// 0 = unsent, 1 = sent, 2 = sent and rcved ack
var (
	ackPack = make(map[int]int)
)

//GBN is the go-back-n algorithm
func GBN(packets []Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr, window int, plp float32) {

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
			fmt.Print("Resending packets ... \n")
			sendWinPack(start, window, packets, conn, addr, noChunks)
		}
	}
}
