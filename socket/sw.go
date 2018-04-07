package socket

import (
	"fmt"
	"net"
	"time"

	"github.com/vmihailenco/msgpack"
)

//SW is the algorithm of stop-and-wait
func SW(packets []Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr, plp float32) {
	for i := 0; i < noChunks; i++ {
		b, err := msgpack.Marshal(&packets[i])
		if err != nil {
			panic(err)
		}

		//drop packets with probability plp
		if i%int(plp*100) != 0 {
			_, err = conn.WriteToUDP(b, addr)
		}
		fmt.Printf("Sent packet %v ... \n", i)

		start := time.Now()
		quit := make(chan int)

		//check if the time exceeded or it received the ack
		go timeAch(start, quit)
		goSend := resendPck(quit)

		if goSend {
			_, err = conn.WriteToUDP(b, addr)
		}
	}
}

// check if time exceeded 0.1 sec
func timeAch(start time.Time, quit chan int) {
	for {
		elapsed := time.Since(start)
		if elapsed > 100000000 {
			quit <- 0
		}
	}
}

// check if we have to resend the packet or not
func resendPck(quit chan int) bool {
	select {
	case <-AckCheck:
		return false
	case <-quit:
		fmt.Println("time exceeded ...")
		return true
	}
}
