package socket

//Packet is the data packet
type Packet struct {
	//Header
	Seqno      uint32
	Cksum, Len uint16
	//Data
	Data []byte
}

//AckPacket is the acknoledgment packet
type AckPacket struct {
	//Header
	Cksum uint16
}
