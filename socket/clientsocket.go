package socket

import (
	"bufio"
	"fmt"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
)

//CreateClientSocket in client-side
func CreateClientSocket(localAddr, servAddr *net.UDPAddr) *net.UDPConn {
	conn, err := net.DialUDP("udp", localAddr, servAddr)
	errors.CheckError(err)
	return conn
}

//SendToServer packets
func SendToServer(conn *net.UDPConn, window int, filename string) {
	p := make([]byte, 2048)
	for i := 0; i < 10; i++ {
		fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")

		_, err := bufio.NewReader(conn).Read(p)
		if err == nil {
			fmt.Printf("%s\n", p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
	}
}
