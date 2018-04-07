package socket

import (
	"fmt"
	"net"

	"github.com/vmihailenco/msgpack"
)

//SW is the algorithm of stop-and-wait
func SW(packets []Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr) {
	for i := 0; i < noChunks; i++ {
		b, err := msgpack.Marshal(&packets[i])
		if err != nil {
			panic(err)
		}
		_, err = conn.WriteToUDP(b, addr)
		fmt.Printf("Sent packet %v ... \n", i)

		// start := time.Now()

		// elapsed := time.Since(start)

		// if elapsed == 0.1 {
		// 	_, err = conn.WriteToUDP(b, addr)
		// }

		//if received ack for the last packet
		if <-AckCheck == uint32(i) {
		}

	}
}