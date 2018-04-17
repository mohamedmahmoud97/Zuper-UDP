package socket

import (
	"net"
)

//SR is the algorithm of selective-repeat
func SR(packets []Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr, window int, plp float32) {
	start := 0
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
				// fmt.Print("Resending packets ... \n")
				// sendWinPack(start, window, packets, conn, addr, noChunks, plp, quit)
			}
		} else {
			ackPack[ackpckt] = 0
			sendWinPack(start, 1, packets, conn, addr, noChunks, plp, quit)
			// b, err := msgpack.Marshal(&packets[pcktseqno])
			// if err != nil {
			// 	panic(err)
			// }
			// _, err = conn.WriteToUDP(b, addr)
		}
	}
}
