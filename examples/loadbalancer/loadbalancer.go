package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/mohamedmahmoud97/Zuper-UDP/errors"
	"github.com/mohamedmahmoud97/Zuper-UDP/loadbalance"
	"github.com/mohamedmahmoud97/Zuper-UDP/socket"
	"github.com/vmihailenco/msgpack"
)

const (
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

func main() {
	//print the logo to the terminal
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER))
		flag.PrintDefaults()
	}
	flag.Parse()
	flag.Usage()

	//Reading Server info from file
	dat, err := ioutil.ReadFile("./device_info/loadbalancer.in")
	errors.CheckError(err)
	servers := strings.Split(string(dat), "\n")

	mainAddr, err := net.ResolveUDPAddr("udp", servers[0])
	errors.CheckError(err)
	mainConn := loadbalance.CreateMainSocket(mainAddr)

	//make the UDP addresses for all servers
	serversAddr := loadbalance.CreateServersAddr(servers[1:])

	serversConn := loadbalance.CreateSockets(serversAddr)

	for {
		buf := make([]byte, 600)
		length, addr, err := mainConn.ReadFromUDP(buf[0:])
		errors.CheckError(err)

		if length > 40 {
			var packet socket.Packet

			err := msgpack.Unmarshal(buf, &packet)
			errors.CheckError(err)

			if loadbalance.IsServer(addr, serversAddr) == true {
				//send the chunks to client that is received from the server
				go loadbalance.SendToClient(mainConn, &packet)
			} else {
				//choose the best server to be assigned to this client request
				bestServer := loadbalance.ChooseServer(serversAddr)
				go loadbalance.AssignToServer(serversConn[bestServer], bestServer, &packet)
			}
		} else if length > 0 && length < 40 {
			var packet socket.AckPacket

			err := msgpack.Unmarshal(buf, &packet)
			errors.CheckError(err)

			if loadbalance.IsServer(addr, serversAddr) == true {
				//send ack to client upon its request
				go loadbalance.SendAckToClient(mainConn, &packet)
			} else {
				//send the ack to the server
				go loadbalance.SendAckToServer(serversConn, &packet)
			}
		}
	}

}
