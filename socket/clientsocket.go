package socket

import (
	"bufio"
	"fmt"
	"net"
)

//CreateClientSocket for any client
func CreateClientSocket(conn *net.UDPConn, window int, filename string) {
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
