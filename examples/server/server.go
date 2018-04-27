package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	server "github.com/mohamedmahmoud97/Zuper-UDP/server"
)

const (
	// BANNER is what is printed for help/info output.
	BANNER = ` 
 _ _ _                          _    _ ___   _ __
|_ _  | _   _ _ __   ___  _ __ | |  | |._ \ | ._ \
   /  /| | | | ._ \ / _ \| '__|| |  | || \ || |_) |
  /  /_| |_| | |_) | |_)/| |   | |__| ||_/ || .__/
 /_ _ _|_ _ _| .__/ \_ _ |_|   |_ __ _|___/ |_|    v1.0  SERVER
             |_|

 A server-side udp reliable data transfer protocol

 #####################################################################

`
)

//ServerInfo is a struct to server info
type ServerInfo struct {
	//PortNumber of the server
	Address string
	//MaxWindow is the max sliding-window size
	MaxWindow int
}

var (
	p, seedValue float32
	//
	lastPort string
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

func getNextSocketAddr(windowSize int) *net.UDPAddr {
	//joining the IP address to the port
	var addr bytes.Buffer
	addr.WriteString(":")
	lastPortInt, _ := strconv.Atoi(lastPort)
	lastPort = strconv.Itoa(lastPortInt + 1)
	addr.WriteString(lastPort)

	serverInfo := ServerInfo{addr.String(), windowSize}
	socketAddr, err := net.ResolveUDPAddr("udp", serverInfo.Address)
	errors.CheckError(err)

	return socketAddr
}

func main() {
	//print the logo to the terminal
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER))
		flag.PrintDefaults()
	}
	flag.Parse()
	flag.Usage()

	//The algorithm to be used in reliability
	algo := os.Args[1]

	//Reading Server info from file
	dat, err := ioutil.ReadFile("./device_info/server.in")
	errors.CheckError(err)
	s := strings.Split(string(dat), "\n")
	port, window, seed, plp := s[0], s[1], s[2], s[3]

	seedValue = convertToFloat(seed)
	p = convertToFloat(plp)
	windowSize, err := strconv.Atoi(window)

	//joining the IP address to the port
	var address bytes.Buffer
	address.WriteString(":")
	address.WriteString(port)

	serverInfo := ServerInfo{address.String(), windowSize}
	servAddr, err := net.ResolveUDPAddr("udp", serverInfo.Address)
	fmt.Printf("started connection in server on port %v ... \n", servAddr)
	errors.CheckError(err)

	//create logfile
	flogS, err := os.OpenFile("serverlog", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	errors.CheckError(err)
	defer flogS.Close()

	//create the socket in server-side
	mainConn := server.CreateSerSocket(servAddr, flogS)

	defer mainConn.Close()

	lastPort = port

	// go read from the connection
	for {
		buf := make([]byte, 600)
		length, addr, err := mainConn.ReadFromUDP(buf[0:])
		errors.CheckError(err)

		if length > 0 {
			socketAddr := getNextSocketAddr(windowSize)
			server.SendAckToClient(mainConn, addr)
			go server.ListenOnSocket(windowSize, algo, p, socketAddr, addr, buf, length)
		}
	}
}
