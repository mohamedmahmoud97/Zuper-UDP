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

	//create the socket between the client and the server
	conn := socket.CreateClientSocket(localAddr, servAddr)
	defer conn.Close()

	//send the filename to the server
	go socket.SendToServer(conn, initWindow, filename)

	//receive any packet from the server
	for {
		// fmt.Println("waiting anything from server ...")
		socket.ReceiveFromServer(conn)
	}
}
