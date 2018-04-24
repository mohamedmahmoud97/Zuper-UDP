package server

import (
	"fmt"
	"log"
	"net"
	"time"

	socket "github.com/mohamedmahmoud97/Zuper-UDP/socket"
	"github.com/vmihailenco/msgpack"
)

var (
	probcountL int
	probcountC = 5
)

//SW is the algorithm of stop-and-wait
func SW(packets []socket.Packet, noChunks int, conn *net.UDPConn, addr *net.UDPAddr, plp float32, AckCheck chan uint32) {
	var ackPack = make(map[int]int)

	for i := 0; i < noChunks; i++ {
		b, err := msgpack.Marshal(&packets[i])
		if err != nil {
			panic(err)
		}

		//drop packets with probability plp
		if probcountL%int(plp*100) != 0 && probcountC%int(plp*100) != 0 {
			_, err = conn.WriteToUDP(b, addr)
			log.SetOutput(flogS)
			log.Printf("Sent packet %v ... \n", i)
			fmt.Printf("Sent packet %v ... \n", i)
		} else {
			log.SetOutput(flogS)
			log.Printf("Sent packet %v but dropped ... \n", i)
			fmt.Printf("Sent packet %v but dropped ... \n", i)
		}

		start := time.Now()
		quit := make(chan uint32)

		//check if the time exceeded or it received the ack
		go timeAch(start, quit, uint32(i), ackPack)
		_, goSend := resendPck(quit, AckCheck)

		if goSend {
			i = i - 1
		} else if !goSend {
			quit <- uint32(i)
		}

		//increment probcount
		probcountL++
		probcountC++
	}
	// log.Printf("Closed connection of port %v... ", conn)
	// fmt.Printf("Closed connection of port %v... ", conn)
	// conn.Close()
}
