package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	socket "github.com/mohamedmahmoud97/Zuper-UDP/socket"
)

//ServerInfo is a struct to server info
type ServerInfo struct {
	//PortNumber of the server
	PortNumber string
	//MaxWindow is the max sliding-window size
	MaxWindow int
}

var (
	p, seedValue float32
)

func convertToFloat(s string) float32 {
	value, err := strconv.ParseFloat(s, 32)
	errors.CheckError(err)
	return float32(value)
}

func convertToInt(s string) int16 {
	value, err := strconv.ParseInt(s, 10, 16)
	errors.CheckError(err)
	return int16(value)
}

func main() {
	//The algorithm to be used in reliability
	algo := os.Args[1]

	//Reading Server info from file
	dat, err := ioutil.ReadFile("/home/mohamedmahmoud/Workspaces/Zuper-UDP/server.in")
	errors.CheckError(err)
	s := strings.Split(string(dat), "\n")
	port, window, seed, ppt := s[0], s[1], s[2], s[3]

	seedValue = convertToFloat(seed)
	p = convertToFloat(ppt)
	windowSize, err := strconv.Atoi(window)

	//joining the IP address to the port
	var address bytes.Buffer
	address.WriteString(":")
	address.WriteString(port)

	server := ServerInfo{address.String(), windowSize}
	servAddr, err := net.ResolveUDPAddr("udp", server.PortNumber)
	fmt.Printf("connection in server on port %v", servAddr)
	errors.CheckError(err)

	//create the socket in server-side
	servConn := socket.CreateSerSocket(servAddr)

	defer servConn.Close()

	// go read from the connection
	for {
		buf := make([]byte, 600)
		length, addr, err := servConn.ReadFromUDP(buf[0:])
		errors.CheckError(err)

		if length > 30 {
			fmt.Print("receiving data packet from clients ... \n")
			go socket.ReceiveReqFromClients(servConn, buf, length, addr, windowSize, algo)
		} else if length > 0 && length < 30 {
			go socket.ReceiveAckFromClients(servConn, buf, length, addr, windowSize, algo)
		}
	}
}
