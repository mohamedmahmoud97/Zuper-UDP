package socket

import (
	"net"
	"time"
)

// a map to give a value for all packets
// 0 = unsent, 1 = sent, 2 = sent and rcved ack
var (
	ackPack  = make(map[int]int)
	pckTimer = make(map[int]time.Time)
	start    int
)

//GBN is the go-back-n algorithm
func GBN(packets []Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr, window int, plp float32) {
	quit := make(chan uint32, window)

	//make all the chunks have value 0 as unack
	for i := 0; i < noChunks; i++ {
		ackPack[i] = 0
	}

	//send the first packets with the window size
	sendWinPack(start, window, packets, conn, addr, noChunks, plp, quit)

	//loop until all the packets are sent and received their ack
	for (start) < noChunks {
		// check if time exceeded or we received a new ack packet
		pcktseqno, goResend := resendPck(quit)
		ackpckt := int(pcktseqno)

		if !goResend {
			if ackpckt == start {
				ackPack[ackpckt] = 2
				start = getNextStart(start, noChunks)
				if start != -1 {
					sendWinPack(start, window, packets, conn, addr, noChunks, plp, quit)
				}
			} else if ackPack[ackpckt] == 1 {
				ackPack[ackpckt] = 2
				nextUnSent := getNextStart(ackpckt, noChunks)
				if nextUnSent < start+4 && nextUnSent != -1 {
					sendWinPack(start, 1, packets, conn, addr, noChunks, plp, quit)
				}
			} else if ackPack[ackpckt] == 2 {

			}
		} else {
			reset(start, ackpckt, window, quit)
			sendWinPack(start, window, packets, conn, addr, noChunks, plp, quit)
		}
	}
	start = 0
}

func reset(start int, ackpckt int, window int, quit chan uint32) {
	for i := ackpckt; i < start+window; i++ {
		ackPack[int(i)] = 0
		quit <- uint32(i)
	}
}
