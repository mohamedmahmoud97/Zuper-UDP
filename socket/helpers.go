package socket

import (
	"fmt"
	"net"

	"github.com/vmihailenco/msgpack"
)

//send a package of packets with the size of the window
//used when start and after changing the start
func sendWinPack(start int, window int, packets []Packet, conn *net.UDPConn, addr *net.UDPAddr, noChunks int) {
	for i := start; i < start+window && i < noChunks; i++ {
		if ackPack[i] == 0 {
			b, err := msgpack.Marshal(&packets[i])
			if err != nil {
				panic(err)
			}
			_, err = conn.WriteToUDP(b, addr)
			ackPack[i] = 1
			fmt.Printf("Sent packet %v ... \n", i)
		}
	}
}

//get the next start which is not sent or not acked
func getNextStart(start int, noChunks int) int {
	for start < noChunks {
		start++
		if ackPack[start] == 0 || ackPack[start] == 1 {
			return start
		}
	}
	//return -1 if it finished sending all packets
	return -1
}
