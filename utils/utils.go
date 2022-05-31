package utils

import (
	"fmt"
	"log"
)

var host string = "http://localhost"

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func PrintConnectionInformation(port string) {
	fmt.Println("Listening on " + host + port)
}
