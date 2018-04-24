package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mohamedmahmoud97/Zuper-UDP/socket"
	"github.com/vmihailenco/msgpack"
)

var (
	probLoss   int
	probCorupt = 5
)

//send a package of packets with the size of the window
//used when start and after changing the start
func sendWinPack(start int, window int, packets []socket.Packet, conn *net.UDPConn, addr *net.UDPAddr, noChunks int, plp float32, quit chan uint32, ackPack map[int]int, pckTimer map[int]time.Time) {
	for i := start; i < start+window && i < noChunks; i++ {
		if ackPack[i] == 0 {
			time.Sleep(1 * time.Millisecond)
			b, err := msgpack.Marshal(&packets[i])
			if err != nil {
				panic(err)
			}

			//drop packets with probability plp
			if probLoss%int(plp*100) != 0 && probCorupt%int(plp*100) != 0 {
				_, err = conn.WriteToUDP(b, addr)
				ackPack[i] = 1

				log.SetOutput(flogS)
				log.Printf("Sent packet %v ... \n", i)
				fmt.Printf("Sent packet %v ... \n", i)
			} else {
				log.SetOutput(flogS)
				log.Printf("Sent packet %v but dropped or corrupted ... \n", i)
				fmt.Printf("Sent packet %v but dropped or corrupted ... \n", i)
			}
			probLoss++
			probCorupt++

			// set timer for each packet
			pckTimer[i] = time.Now()
			go timeAch(pckTimer[i], quit, uint32(i), ackPack)
		}
	}
}

//get the next start which is not sent or not acked
func getNextStart(start int, noChunks int, ackPack map[int]int) int {
	for start < noChunks {
		start++
		if ackPack[start] == 0 || ackPack[start] == 1 {
			return start
		}
	}
	//return -1 if it finished sending all packets
	return -1
}

// check if time exceeded 0.1 sec
func timeAch(start time.Time, quit chan uint32, seqno uint32, ackPack map[int]int) {
	for {
		select {
		case <-quit:
			return
		default:
			elapsed := time.Since(start)
			if elapsed > 200000000 {
				if ackPack[int(seqno)] != 2 {
					log.SetOutput(flogS)
					log.Printf("time exceeded for pckt %v\n", seqno)
					fmt.Printf("time exceeded for pckt %v\n", seqno)
					quit <- seqno
				}
				return
			}
		}
	}
}

// check if we have to resend the packet or not
func resendPck(quit chan uint32, AckCheck chan uint32) (uint32, bool) {
	select {
	case ackseqno := <-AckCheck:
		return ackseqno, false
	case timeoutpck := <-quit:
		return timeoutpck, true
	}
}
