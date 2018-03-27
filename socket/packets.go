package socket

//Packet is the data packet
type Packet struct {
	//Header
	cksum, len uint16
	seqno      uint32
	//Data
	data []byte
}

//AckPacket is the acknoledgment packet
type AckPacket struct {
	//Header
	cksum uint16
}
