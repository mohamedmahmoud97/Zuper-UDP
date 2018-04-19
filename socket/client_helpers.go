package socket

import (
	"fmt"
	"time"
)

// check if time exceeded 0.1 sec
func fileTimer(start time.Time, quit chan uint32) {
	for {
		select {
		case <-quit:
			return
		default:
			elapsed := time.Since(start)
			if elapsed > 100000000 {
				fmt.Printf("time exceeded for file request pckt ack ...\n")
				quit <- 1
				return
			}
		}
	}
}

// check if we have to resend the packet or not
func resendReq(quit chan uint32) bool {
	select {
	case <-AckFileCheck:
		fmt.Println("Received an ack from the server ... ")
		return false
	case <-quit:
		fmt.Println("Will resend the requested file again ... ")
		return true
	}
}
