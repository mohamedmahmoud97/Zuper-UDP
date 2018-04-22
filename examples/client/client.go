package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	client "github.com/mohamedmahmoud97/Zuper-UDP/client"
	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
)

func convertToFloat(s string) float32 {
	value, err := strconv.ParseFloat(s, 32)
	errors.CheckError(err)
	return float32(value)
}

func main() {
	//The algorithm to be used in reliability
	algo := os.Args[1]

	//Reading Server info from file
	dat, err := ioutil.ReadFile("../../device_info/client.in")
	errors.CheckError(err)
	s := strings.Split(string(dat), "\n")
	servIP, servPort, clientIP, clientPort, filename, window, prob := s[0], s[1], s[2], s[3], s[4], s[5], s[6]

	initWindow, err := strconv.Atoi(window)
	plp := convertToFloat(prob)

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
	localAddr, err := net.ResolveUDPAddr("udp", cAddress.String())
	errors.CheckError(err)

	//create the socket between the client and the server
	conn := client.CreateClientSocket(localAddr, servAddr)
	defer conn.Close()

	//create logfile
	flogC, err := os.OpenFile("clientlog", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	errors.CheckError(err)
	defer flogC.Close()

	//send the filename to the server
	go client.SendToServer(conn, initWindow, filename, plp, flogC)

	// go read from the connection
	for {
		buf := make([]byte, 600)
		length, addr, err := conn.ReadFromUDP(buf[0:])
		errors.CheckError(err)

		if length > 25 {
			go client.ReceiveFromServer(conn, buf, addr, algo)
		} else if length > 0 && length < 25 {
			go client.ReceiveAckFromServer()
		}
	}
}
