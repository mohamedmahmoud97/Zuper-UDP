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

	"github.com/mohamedmahmoud97/Zuper-UDP/errors"
	"github.com/mohamedmahmoud97/Zuper-UDP/loadbalance"
	"github.com/mohamedmahmoud97/Zuper-UDP/socket"
	"github.com/vmihailenco/msgpack"
)

const (
	p = 0.1
	// BANNER is what is printed for help/info output.
	BANNER = ` 
 _ _ _                          _    _ ___   _ __
|_ _  | _   _ _ __   ___  _ __ | |  | |._ \ | ._ \
   /  /| | | | ._ \ / _ \| '__|| |  | || \ || |_) |
  /  /_| |_| | |_) | |_)/| |   | |__| ||_/ || .__/
 /_ _ _|_ _ _| .__/ \_ _ |_|   |_ __ _|___/ |_|    v1.0  LOADBALANCER
             |_|

 A server-side loadbalancer with udp reliable data transfer protocol

 #####################################################################

`
)

var (
	lastPort string
	buffer   = make(map[*net.UDPAddr]map[int][]byte)
)

func getNextSocketAddr() string {
	//joining the IP address to the port
	var addr bytes.Buffer
	lastPortInt, _ := strconv.Atoi(lastPort)
	lastPort = strconv.Itoa(lastPortInt + 1)
	addr.WriteString(lastPort)
	return addr.String()
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

	//Reading Loadbalancer info from file
	dat, err := ioutil.ReadFile("./device_info/loadbalancer.in")
	errors.CheckError(err)
	servers := strings.Split(string(dat), "\n")

	mainAddr, err := net.ResolveUDPAddr("udp", servers[0])
	errors.CheckError(err)

	//create logfile
	flogL, err := os.OpenFile("loadbalancerlog", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	errors.CheckError(err)
	defer flogL.Close()

	//create main socket
	mainConn := loadbalance.CreateMainSocket(mainAddr, flogL)

	//make the UDP addresses for all servers
	serversAddr := loadbalance.CreateServersAddr(servers)
	// serverIP := serversAddr[0].IP
	serverPort := serversAddr[0].Port
	lastPort = strconv.Itoa(serverPort)

	for {
		buf := make([]byte, 700)
		length, addr, err := mainConn.ReadFromUDP(buf[0:])
		errors.CheckError(err)

		if length > 115 {
			// receiving data packets
			var packet socket.Packet

			err := msgpack.Unmarshal(buf, &packet)
			errors.CheckError(err)

			if loadbalance.IsServer(addr, serversAddr) == true {
				//send the chunks to client that is received from the server
				go loadbalance.SendToClient(mainConn, &packet)
				go loadbalance.AppendToBuffer(&packet, algo)
			} else {
				// path := "./" + string(packet.Data)
				// if _, err := os.Stat(path); err == nil {
				// 	socketPort := getNextSocketAddr()

				// 	//joining the IP address to the port
				// 	var bind bytes.Buffer
				// 	bind.WriteString(serverIP.String())
				// 	bind.WriteString(":")
				// 	bind.WriteString(socketPort)

				// 	socketAddr, err := net.ResolveUDPAddr("udp", bind.String())
				// 	errors.CheckError(err)

				// 	go server.ListenOnSocket(5, algo, p, socketAddr, addr, &packet)
				// } else {
				// 	//add filename to buffer
				// 	go loadbalance.SaveFileName(addr, string(packet.Data))

				// }
				//choose the best server to be assigned to this client request
				bestServer := loadbalance.ChooseServer(serversAddr)
				go loadbalance.AssignToServer(mainConn, bestServer, &packet)
			}
		} else if length > 0 && length < 115 {
			// receiving ack packets
			var packet socket.AckPacket

			err := msgpack.Unmarshal(buf, &packet)
			errors.CheckError(err)

			if loadbalance.IsServer(addr, serversAddr) == true {
				//send ack to client upon its request
				go loadbalance.SendAckToClient(mainConn, addr, &packet)
			} else {
				//send the ack to the server
				go loadbalance.SendAckToServer(mainConn, &packet)
			}
		}
	}

}
