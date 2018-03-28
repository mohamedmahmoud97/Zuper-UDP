package socket

//Packet is the data packet
type Packet struct {
	//Header
	Seqno      uint32
	Cksum, Len uint16
	//Data
	Data []byte
}

//CheckSum of the packet
func (p *Packet) CheckSum() {
	//cksum := 0

}

//AckPacket is the acknoledgment packet
type AckPacket struct {
	//Header
	Seqno uint32
	Cksum uint16
}
