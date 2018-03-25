package main

import (
	"net"
)

const (
	//ServerIP of the server
	ServerIP = "192.168.1.11"
	//ServerPort number
	ServerPort = "5000"
	//ClientPort number
	ClientPort = "3000"
)

var (
	//FileName to be transferred
	FileName string
	//WindowSize is the initial sliding-window size
	WindowSize int
)

func main() {
	// initialize all connections
	servAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	CheckError(err)
	listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10002")
	CheckError(err)
	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)
	conn, err := net.DialUDP("udp", localAddr, servAddr)
	CheckError(err)
	servConn, err := net.ListenUDP("udp", listenAddr)
	CheckError(err)
}
