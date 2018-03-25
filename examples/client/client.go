package main

import (
	"bufio"
	"fmt"
	"net"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
)

const (
	//ServerIP of the server
	serverAddr = "127.0.0.1:10001"
	//ServerPort number
	serverPort = "10001"
	//ClientPort number
	clientAddr = "127.0.0.1:10000"
)

var (
	//FileName to be transferred
	FileName string
	//WindowSize is the initial sliding-window size
	WindowSize int
)

func main() {
	// initialize all connections
	servAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	errors.CheckError(err)
	// /listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10002")
	errors.CheckError(err)
	localAddr, err := net.ResolveUDPAddr("udp", clientAddr)
	errors.CheckError(err)
	conn, err := net.DialUDP("udp", localAddr, servAddr)
	errors.CheckError(err)
	//servConn, err := net.ListenUDP("udp", listenAddr)
	errors.CheckError(err)

	p := make([]byte, 2048)
	for {
		fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")

		_, err = bufio.NewReader(conn).Read(p)
		if err == nil {
			fmt.Printf("%s\n", p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
	}
	conn.Close()

	// create a channel for a packet number to be written to
	//i := make(chan int, 1)
	//go func () {
	//loop:
	//// wait for the ack while we're waiting for a packet or timing out
	//	go wait_for_ack(servConn, packetnum, i)
	//	for {
	//
	//		select {
	//		case res := <-i:
	//			fmt.Println("\nPacket accepted!")
	//			packet_num = res+1
	//			// wait for another ack for the next one if we get the right packet
	//			goto loop
	//		case <-time.After(100 * time.Millisecond):
	//			fmt.Println("timed out for", packet_num)
	//			// if it takes too long for an ACK, go send the packet again
	//			write(conn)
	//		}
	//	}
	//}()
	//// go write to the connection because the previous stuff is
	//// all hanging out in the background for now
	//for {
	//	write(conn)
	//}
	defer conn.Close()
}
