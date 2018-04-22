package socket

//Packet is the data packet
type Packet struct {
	//Header
	Seqno      uint32
	PckNo, Len uint16
	//Data
	Data []byte
}

//AckPacket is the acknoledgment packet
type AckPacket struct {
	//Header
	Seqno uint32
	Cksum uint16
}

//CheckSum of the packet
func (p *Packet) CheckSum() {

}
