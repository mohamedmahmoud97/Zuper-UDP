package client

import (
	"fmt"
	"hash/adler32"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sort"
	"time"

	errors "github.com/mohamedmahmoud97/Zuper-UDP/errors"
	socket "github.com/mohamedmahmoud97/Zuper-UDP/socket"
	"github.com/vmihailenco/msgpack"
)

var (
	//AckFileCheck is channel on receiving ack on file request
	AckFileCheck       = make(chan uint32)
	lastAck      int32 = -1
	buffer             = make(map[int][]byte)
	corruptProb  int
	fileName     string
	pckNo        uint16
	plp          float32
	flogC        *os.File
)

//CreateClientSocket in client-side
func CreateClientSocket(localAddr *net.UDPAddr) *net.UDPConn {
	conn, err := net.ListenUDP("udp", localAddr)
	errors.CheckError(err)
	return conn
}

//SendToServer the filename of the needed file
func SendToServer(conn *net.UDPConn, servAddr, localAddr *net.UDPAddr, window int, filename string, flogc *os.File) {
	flogC = flogc
	log.SetOutput(flogC)
	log.Printf("client is requesting file %v from server ... \n", filename)

	fmt.Printf("client is requesting file %v from server ... \n", filename)
	fileName = filename

	file := []byte(filename)

	noOfBytes := uint16(len(file))
	reqPacket := socket.Packet{Data: file, PckNo: 1, Len: noOfBytes, SrcAddr: localAddr, DstAddr: servAddr}

	b, err := msgpack.Marshal(&reqPacket)
	if err != nil {
		panic(err)
	}

	log.SetOutput(flogC)
	log.Println("Encoded the message ...")
	fmt.Println("Encoded the message ...")

	//send the message to the server
	_, err = conn.WriteToUDP(b, servAddr)
	errors.CheckError(err)

	// conn.Close()

	start := time.Now()
	quit := make(chan uint32)

	//check if the time exceeded or it received the ack
	go fileTimer(start, quit)
	exists, goSend := resendReq(quit)

	if goSend {
		SendToServer(conn, servAddr, localAddr, window, filename, flogC)
	} else if !goSend {
		if exists == 0 {
			//file doesn't exists terminate program
			log.Println("File doesn't exists ...")
			fmt.Println("File doesn't exists ...")
			os.Exit(1)
		} else if exists == 1 {
			//file exists don't do anything
		}
		quit <- 0
	}
}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr, packet *socket.Packet) {
	ack := socket.AckPacket{Seqno: packet.Seqno, SrcAddr: packet.DstAddr, DstAddr: packet.SrcAddr}

	b, err := msgpack.Marshal(&ack)
	if err != nil {
		panic(err)
	}

	_, err = conn.WriteToUDP(b, addr)
	errors.CheckError(err)
}

//ReceiveFromServer any ack packet
func ReceiveFromServer(conn *net.UDPConn, buf []byte, addr *net.UDPAddr, algo string) {
	var packet socket.Packet

	err := msgpack.Unmarshal(buf, &packet)
	if err != nil {
		panic(err)
	}

	log.SetOutput(flogC)
	log.Printf("Delivered packet with seqno %v \n", packet.Seqno)
	fmt.Printf("Delivered packet with seqno %v \n", packet.Seqno)

	if packet.Cksum == adler32.Checksum(packet.Data) {
		if algo == "sw" {
			go sendResponse(conn, addr, &packet)
			appendFile(packet.Data, packet.Seqno)
			done := CheckOnPck(&packet, algo, buffer, fileName)
			if done == 1 {
				os.Exit(0)
			}
		} else if algo == "gbn" {
			if int32(packet.Seqno) == lastAck+1 {
				lastAck = int32(packet.Seqno)
				fmt.Printf("last ack packet is %v\n", lastAck)
				appendFile(packet.Data, packet.Seqno)
				go sendResponse(conn, addr, &packet)
				done := CheckOnPck(&packet, algo, buffer, fileName)
				if done == 1 {
					os.Exit(0)
				}
			} else if int32(packet.Seqno) > lastAck+1 && lastAck != -1 {
				//change seqno of ack packet to last delivered packet
				packet.Seqno = uint32(lastAck)
				go sendResponse(conn, addr, &packet)
			} else if int32(packet.Seqno) > lastAck+1 && lastAck == -1 {

			}
		} else if algo == "sr" {
			appendFile(packet.Data, packet.Seqno)
			go sendResponse(conn, addr, &packet)
			done := CheckOnPck(&packet, algo, buffer, fileName)
			if done == 1 {
				os.Exit(0)
			}
		}
	}
}

//ReceiveAckFromServer any packet
func ReceiveAckFromServer(buf []byte) {
	var packet socket.AckPacket

	err := msgpack.Unmarshal(buf, &packet)
	if err != nil {
		panic(err)
	}

	log.SetOutput(flogC)
	log.Println("Received Ack of requested file packet ...")
	fmt.Printf("Received Ack of requested file packet ... \n")

	//a channel for sending seqno
	AckFileCheck <- packet.Seqno
}

///append to buffer to build file later on
func appendFile(data []byte, seqno uint32) {
	buffer[int(seqno)] = data
}

//build the requested file at the client-side
func buildFile(algo string, buffer map[int][]byte, filename string) {
	log.SetOutput(flogC)
	log.Println("Building File ... ")
	fmt.Println("Building File ... ")

	recData := make([]byte, pckNo*512)

	if algo != "sw" {
		// To store the keys in slice in sorted order
		var keys []int
		for k := range buffer {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		// to store sorted buffer
		for _, k := range keys {
			for i := 0; i < len(buffer[k]); i++ {
				recData = append(recData, buffer[k][i])
			}
		}
	} else {
		// for sw algorithm no sorting just putting buffer to file
		for k := range buffer {
			for i := 0; i < len(buffer[k]); i++ {
				recData = append(recData, buffer[k][i])
			}
		}
	}

	err := ioutil.WriteFile(filename, recData, 0644)
	errors.CheckError(err)
	log.SetOutput(flogC)
	log.Println("Finished ... ")
	fmt.Println("Finished ... ")
}

//CheckOnPck is checking if the packet is the last packet or not
func CheckOnPck(packet *socket.Packet, algo string, buffer map[int][]byte, filename string) int {
	if packet.Seqno == 0 {
		pckNo = packet.PckNo
		return 0
	} else if int(packet.Seqno) == int(pckNo)-1 {
		fmt.Println(filename)
		time.Sleep(5 * time.Millisecond)
		buildFile(algo, buffer, filename)
		return 1
	} else {
		return 0
	}
}
