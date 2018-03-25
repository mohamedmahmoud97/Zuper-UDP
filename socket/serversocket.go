import(
	"net"
	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
)

var(
	//Data
	Data make([][512]byte)
)

def sendResponse(conn *net.UDPConn, addr net.UDPAddr){
	_, err := conn.WriteToUDP([]byte("An ack to the sent packet from client"), addr)
	errors.CheckError(err)
}

def ReceiveFromClients(conn *net.UDPConn){
	var buf [512]byte

    _, addr, err := conn.ReadFromUDP(buf[0:])
    if err != nil {
        return
	}
	
	// Data.append(buf)

	go sendResponse(conn, addr)
}