package socket

//Packet is the data packet
type Packet struct {
	//Header
	srcPort, dstPort string
	cksum, len       uint16
	seqno            uint32
	//Data
	data [500]byte
}

//AckPacket is the acknoledgment packet
type AckPacket struct {
	//Header
	srcPort, dstPort string
	cksum            uint16
}
