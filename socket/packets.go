package socket

import (
	"net"
)

//Packet is the data packet
type Packet struct {
	//Header
	Seqno      uint32
	PckNo, Len uint16
	//Data
	Data []byte
	//Cksum
	Cksum uint32
	//Addresses
	SrcAddr, DstAddr *net.UDPAddr
}

//AckPacket is the acknoledgment packet
type AckPacket struct {
	//Header
	Seqno uint32
	Cksum uint32
	//Addresses
	SrcAddr, DstAddr *net.UDPAddr
}
