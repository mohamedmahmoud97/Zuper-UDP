package socket

import (
	"net"

	"github.com/vmihailenco/msgpack"
)

//GBN is the go-back-n algorithm
func GBN(packets []Packet, noChunks int, conn *net.UDPConn, servAddr *net.UDPAddr) {
	for i := 0; i < noChunks; i++ {
		b, err := msgpack.Marshal(&packets[i])
		if err != nil {
			panic(err)
		}
		_, err = conn.Write(b)
	}
}
