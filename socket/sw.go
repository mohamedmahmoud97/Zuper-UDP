package socket

import (
	"fmt"
	"net"
	"time"

	"github.com/vmihailenco/msgpack"
)

var (
	probcount int
)

//SW is the algorithm of stop-and-wait
func SW(packets []Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr, plp float32) {
	for i := 0; i < noChunks; i++ {
		b, err := msgpack.Marshal(&packets[i])
		if err != nil {
			panic(err)
		}

		//drop packets with probability plp
		if probcount%int(plp*100) != 0 {
			_, err = conn.WriteToUDP(b, addr)
			fmt.Printf("Sent packet %v ... \n", i)
		} else {
			fmt.Printf("Sent packet %v but dropped ... \n", i)
		}

		start := time.Now()
		quit := make(chan uint32)

		//check if the time exceeded or it received the ack
		go timeAch(start, quit, uint32(i))
		_, goSend := resendPck(quit)

		if goSend {
			i = i - 1
		} else if !goSend {
			quit <- uint32(i)
		}

		//increment probcount
		probcount++
	}
}
