package socket

import (
	"net"

	"github.com/vmihailenco/msgpack"
)

//SR is the algorithm of selective-repeat
func SR(packets []Packet, noChunks int, conn *net.UDPConn, servAddr *net.UDPAddr, window int) {
	for i := 0; i < noChunks; i++ {
		b, err := msgpack.Marshal(&packets[i])
		if err != nil {
			panic(err)
		}
		_, err = conn.Write(b)
	}
}
