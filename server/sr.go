package server

import (
	"net"
	"time"

	socket "github.com/mohamedmahmoud97/Zuper-UDP/socket"
)

//SR is the algorithm of selective-repeat
func SR(packets []socket.Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr, window int, plp float32, AckCheck chan uint32) {
	var ackPack = make(map[int]int)
	var pckTimer = make(map[int]time.Time)
	start := 0
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
		pcktseqno, goResend := resendPck(quit, AckCheck)
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
				if nextUnSent < start+4 && nextUnSent != -1 {
					sendWinPack(start, 1, packets, conn, addr, noChunks, plp, quit, ackPack, pckTimer)
				}
			} else if ackPack[ackpckt] == 2 {

			}
		} else {
			if ackPack[ackpckt] != 2 {
				ackPack[ackpckt] = 0
				// time.Sleep(1 * time.Millisecond)
				sendWinPack(start, 1, packets, conn, addr, noChunks, plp, quit, ackPack, pckTimer)
			}
		}
	}
	start = 0
}
