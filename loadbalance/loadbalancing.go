package loadbalance

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mohamedmahmoud97/Zuper-UDP/errors"
	"github.com/mohamedmahmoud97/Zuper-UDP/socket"
	"github.com/vmihailenco/msgpack"
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

//IsServer is to check if the incoming packet is from a server or not
func IsServer(addr *net.UDPAddr, servers []*net.UDPAddr) bool {
	for i := 0; i < len(servers); i++ {
		if addr == servers[i] {
			return true
		}
	}
	return false
}

//AssignToServer a request from client
func AssignToServer(conn *net.UDPConn, server *net.UDPAddr, packet *socket.Packet) {
	pck := socket.Packet{Data: packet.Data, PckNo: packet.PckNo, Len: packet.Len, SrcAddr: packet.SrcAddr, DstAddr: server}
	b, err := msgpack.Marshal(&pck)
	if err != nil {
		panic(err)
	}

	//send the message to the server
	_, err = conn.WriteToUDP(b, server)
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
func SendAckToClient(mainConn *net.UDPConn, packet *socket.Packet) {
	ack := socket.AckPacket{Seqno: packet.Seqno, SrcAddr: packet.SrcAddr, DstAddr: packet.DstAddr}

	b, err := msgpack.Marshal(&ack)
	if err != nil {
		panic(err)
	}

	_, err = mainConn.WriteToUDP(b, packet.DstAddr)
	errors.CheckError(err)

	log.SetOutput(flogL)
	log.Printf("SenT ack packet from server %v to client %v ...\n", packet.SrcAddr, packet.DstAddr)
	fmt.Printf("SenT ack packet from server %v to client %v ...\n", packet.SrcAddr, packet.DstAddr)
}

//SendAckToServer which is received from the client
func SendAckToServer(connections map[*net.UDPAddr]*net.UDPConn, packet *socket.Packet) {
	ack := socket.AckPacket{Seqno: packet.Seqno, SrcAddr: packet.SrcAddr, DstAddr: packet.DstAddr}

	b, err := msgpack.Marshal(&ack)
	if err != nil {
		panic(err)
	}

	serverConn := connections[packet.DstAddr]

	_, err = serverConn.WriteToUDP(b, packet.DstAddr)
	errors.CheckError(err)

	log.SetOutput(flogL)
	log.Printf("SenT ack packet %v from client %v to server %v ...\n", packet.Seqno, packet.SrcAddr, packet.DstAddr)
	fmt.Printf("SenT ack packet %v from client %v to server %v ...\n", packet.Seqno, packet.SrcAddr, packet.DstAddr)
}
