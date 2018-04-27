package loadbalance

import (
	"net"
	"os"

	"github.com/mohamedmahmoud97/Zuper-UDP/errors"
)

var (
	flogL *os.File
	n     int
)

//CreateMainSocket for loadbalancer to listen on this port
func CreateMainSocket(addr *net.UDPAddr) *net.UDPConn {
	//create the socket on the port number
	servConn, err := net.ListenUDP("udp", addr)
	errors.CheckError(err)
	return servConn
}

//CreateServersAddr is creating UDP addresses for all servers
func CreateServersAddr(servers []string) []*net.UDPAddr {
	serversAddr := make([]*net.UDPAddr, len(servers))
	for i := 0; i < len(servers); i++ {
		servAddr, err := net.ResolveUDPAddr("udp", servers[i])
		errors.CheckError(err)
		serversAddr = append(serversAddr, servAddr)
	}
	return serversAddr
}

//CreateSockets between all servers and the loadbalancer
func CreateSockets(servers []*net.UDPAddr) map[*net.UDPAddr]*net.UDPConn {
	//a map of addresses with its connections
	serversConn := make(map[*net.UDPAddr]*net.UDPConn, len(servers))
	for i := 0; i < len(servers); i++ {
		//create the socket on the port number
		servConn, err := net.ListenUDP("udp", servers[i])
		errors.CheckError(err)
		serversConn[servers[i]] = servConn
	}
	return serversConn
}

//ChooseServer with weightedRoundRobin algorithm
func ChooseServer(serversAddr []*net.UDPAddr) *net.UDPAddr {
	temp := n % len(serversAddr)
	n++
	return serversAddr[temp]
}

//AssignToServer a request from client
func AssignToServer(conn *net.UDPConn, buf []byte) {

}

//IsServer is to check if the incoming packet is from a server or not
func IsServer(addr *net.UDPAddr, servers []*net.UDPAddr) bool {
	for i := 0; i < len(servers); i++ {
		if addr == servers[i] {
			return true
		}
	}
	return false
}

//SendToClient the datagrams received from the server
func SendToClient() {

}

//SendAckToClient which received from the server
func SendAckToClient() {

}

//SendAckToServer which is received from the client
func SendAckToServer() {

}
