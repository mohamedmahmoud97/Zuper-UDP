package socket

import (
	"fmt"
	"net"

	"github.com/vmihailenco/msgpack"
)

//SW is the algorithm of stop-and-wait
func SW(packets []Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr) {
	// ackChan := make(chan int)

	for i := 0; i < noChunks; i++ {
		b, err := msgpack.Marshal(&packets[i])
		if err != nil {
			panic(err)
		}
		_, err = conn.WriteToUDP(b, addr)
		fmt.Printf("Sent packet %v ... \n", i)

		// ackCheck := 0

		// ackCheck <- ackChan
	}
}
