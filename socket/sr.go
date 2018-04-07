package socket

import (
	"net"
)

//SR is the algorithm of selective-repeat
func SR(packets []Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr, window int) {
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
}
