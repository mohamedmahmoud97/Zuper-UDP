package server

import (
	"net"
	"time"

	socket "github.com/mohamedmahmoud97/Zuper-UDP/socket"
)

//GBN is the go-back-n algorithm
func GBN(packets []socket.Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr, window int, plp float32) {
	// a map to give a value for all packets
	// 0 = unsent, 1 = sent, 2 = sent and rcved ack
	var ackPack = make(map[int]int)
	var pckTimer = make(map[int]time.Time)
	var start int
	quit := make(chan uint32, window)

	//make all the chunks have value 0 as unack
	for i := 0; i < noChunks; i++ {
		ackPack[i] = 0
	}

	//send the first packets with the window size
	sendWinPack(start, window, packets, conn, addr, noChunks, plp, quit, ackPack, pckTimer)

	//loop until all the packets are sent and received their ack
	for (start) < noChunks {
		// check if time exceeded or we received a new ack packet
		pcktseqno, goResend := resendPck(quit)
		ackpckt := int(pcktseqno)

		if !goResend {
			if ackpckt == start {
				ackPack[ackpckt] = 2
				time.Sleep(1 * time.Millisecond)
				start = getNextStart(start, noChunks, ackPack)
				if start != -1 {
					sendWinPack(start, window, packets, conn, addr, noChunks, plp, quit, ackPack, pckTimer)
				}
			} else if ackPack[ackpckt] == 1 {
				ackPack[ackpckt] = 2
				time.Sleep(1 * time.Millisecond)
				nextUnSent := getNextStart(ackpckt, noChunks, ackPack)
				if nextUnSent < start+window && nextUnSent != -1 {
					sendWinPack(start, 1, packets, conn, addr, noChunks, plp, quit, ackPack, pckTimer)
				}
			} else if ackPack[ackpckt] == 2 {
				time.Sleep(1 * time.Millisecond)
			}
		} else {
			reset(start, ackpckt, window, quit, ackPack)
			sendWinPack(start, window, packets, conn, addr, noChunks, plp, quit, ackPack, pckTimer)
		}
	}
	start = 0
}

func reset(start int, ackpckt int, window int, quit chan uint32, ackPack map[int]int) {
	for i := start; i < start+window; i++ {
		if ackPack[i] != 2 {
			ackPack[i] = 0
			time.Sleep(1 * time.Millisecond)
		}
		quit <- uint32(i)
	}
}
