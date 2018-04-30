package loadbalance

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/mohamedmahmoud97/Zuper-UDP/errors"
	"github.com/mohamedmahmoud97/Zuper-UDP/socket"
	"github.com/vmihailenco/msgpack"
)

var (
	flogL *os.File
	n     int
)

//CreateMainSocket for loadbalancer to listen on this port
func CreateMainSocket(addr *net.UDPAddr, log *os.File) *net.UDPConn {
	flogL = log

	//create the socket on the port number
	servConn, err := net.ListenUDP("udp", addr)
	errors.CheckError(err)
	return servConn
}

//CreateServersAddr is creating UDP addresses for all servers
func CreateServersAddr(servers []string) []*net.UDPAddr {
	serversAddr := make([]*net.UDPAddr, len(servers)-1)
	for i, addr := range servers[1:] {
		servAddr, err := net.ResolveUDPAddr("udp", addr)
		errors.CheckError(err)
		serversAddr[i] = servAddr
	}
	return serversAddr
}

//ChooseServer with weightedRoundRobin algorithm
func ChooseServer(serversAddr []*net.UDPAddr) *net.UDPAddr {
	temp := n % len(serversAddr)
	n++
	return serversAddr[temp]
}

//IsServer is to check if the incoming packet is from a server or not
func IsServer(addr *net.UDPAddr, servers []*net.UDPAddr) bool {
	str := addr.String()
	strs := strings.Split(str, ":")
	ip := strs[0]

	for i := 0; i < len(servers); i++ {
		strS := servers[i].String()
		strsS := strings.Split(strS, ":")
		ipS := strsS[0]
		if ip == ipS {
			return true
		}
	}
	return false
}

//AssignToServer a request from client
func AssignToServer(mainConn *net.UDPConn, server *net.UDPAddr, packet *socket.Packet) {
	pck := socket.Packet{Data: packet.Data, PckNo: packet.PckNo, Len: packet.Len, SrcAddr: packet.SrcAddr, DstAddr: server}
	b, err := msgpack.Marshal(&pck)
	if err != nil {
		panic(err)
	}

	//send the message to the server
	_, err = mainConn.WriteToUDP(b, server)
	errors.CheckError(err)

	log.SetOutput(flogL)
	log.Printf("Assigned server %v to client %v ...\n", server, packet.SrcAddr)
	fmt.Printf("Assigned server %v to client %v ...\n", server, packet.SrcAddr)
}

//SendToClient the datagrams received from the server
func SendToClient(mainConn *net.UDPConn, packet *socket.Packet) {
	pck := socket.Packet{Data: packet.Data, Len: packet.Len, Seqno: packet.Seqno, PckNo: packet.PckNo, Cksum: packet.Cksum,
		SrcAddr: packet.SrcAddr, DstAddr: packet.DstAddr}
	b, err := msgpack.Marshal(&pck)
	if err != nil {
		panic(err)
	}

	//send the message to the server
	_, err = mainConn.WriteToUDP(b, packet.DstAddr)
	errors.CheckError(err)

	log.SetOutput(flogL)
	log.Printf("SenT packet %v from server %v to client %v ...\n", packet.Seqno, packet.SrcAddr, packet.DstAddr)
	fmt.Printf("SenT packet %v from server %v to client %v ...\n", packet.Seqno, packet.SrcAddr, packet.DstAddr)

}

//SendAckToClient which received from the server
func SendAckToClient(mainConn *net.UDPConn, addr *net.UDPAddr, packet *socket.AckPacket) {
	ack := socket.AckPacket{Seqno: packet.Seqno, SrcAddr: addr, DstAddr: packet.DstAddr}

	b, err := msgpack.Marshal(&ack)
	if err != nil {
		panic(err)
	}

	_, err = mainConn.WriteToUDP(b, packet.DstAddr)
	errors.CheckError(err)

	log.SetOutput(flogL)
	log.Printf("SenT ack packet from server %v to client %v ...\n", addr, packet.DstAddr)
	fmt.Printf("SenT ack packet from server %v to client %v ...\n", addr, packet.DstAddr)
}

//SendAckToServer which is received from the client
func SendAckToServer(mainConn *net.UDPConn, packet *socket.AckPacket) {
	ack := socket.AckPacket{Seqno: packet.Seqno, SrcAddr: packet.SrcAddr, DstAddr: packet.DstAddr}

	b, err := msgpack.Marshal(&ack)
	if err != nil {
		panic(err)
	}

	_, err = mainConn.WriteToUDP(b, packet.DstAddr)
	errors.CheckError(err)

	log.SetOutput(flogL)
	log.Printf("SenT ack packet %v from client %v to server %v ...\n", packet.Seqno, packet.SrcAddr, packet.DstAddr)
	fmt.Printf("SenT ack packet %v from client %v to server %v ...\n", packet.Seqno, packet.SrcAddr, packet.DstAddr)
}
