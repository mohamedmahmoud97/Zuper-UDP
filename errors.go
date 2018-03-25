package errors

import (
	"fmt"
	"os"
)

//CheckError is good for checking errors
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}
