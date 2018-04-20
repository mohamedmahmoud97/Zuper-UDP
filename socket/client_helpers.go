package socket

import (
	"fmt"
	"log"
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
				log.SetOutput(flogC)
				log.Println("time exceeded for file request pckt ack ...")
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
		log.SetOutput(flogC)
		log.Println("Received an ack from the server ... ")
		fmt.Println("Received an ack from the server ... ")
		return false
	case <-quit:
		log.SetOutput(flogC)
		log.Println("Will resend the requested file again ... ")
		fmt.Println("Will resend the requested file again ... ")
		return true
	}
}
