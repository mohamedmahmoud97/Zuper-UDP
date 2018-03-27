package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	socket "github.com/mohamedmahmoud97/Zuper-UDP/socket"
)

func main() {
	//Reading Server info from file
	dat, err := ioutil.ReadFile("/home/mohamedmahmoud/Workspaces/Zuper-UDP/client.in")
	errors.CheckError(err)
	s := strings.Split(string(dat), "\n")
	servIP, servPort, clientIP, clientPort, filename, window := s[0], s[1], s[2], s[3], s[4], s[5]

	initWindow, err := strconv.Atoi(window)

	fmt.Println(initWindow, filename)

	//joining the IP address to the port
	var sAddress bytes.Buffer
	sAddress.WriteString(servIP)
	sAddress.WriteString(":")
	sAddress.WriteString(servPort)

	var cAddress bytes.Buffer
	cAddress.WriteString(clientIP)
	cAddress.WriteString(":")
	cAddress.WriteString(clientPort)

	// initialize all connections
	servAddr, err := net.ResolveUDPAddr("udp", sAddress.String())
	errors.CheckError(err)
	//listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10002")
	errors.CheckError(err)
	localAddr, err := net.ResolveUDPAddr("udp", cAddress.String())
	errors.CheckError(err)

	//servConn, err := net.ListenUDP("udp", listenAddr)
	errors.CheckError(err)

	conn := socket.CreateClientSocket(localAddr, servAddr)
	defer conn.Close()

	socket.SendToServer(conn, initWindow, filename)

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
}
