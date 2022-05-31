package utils

import (
	"crypto/sha256"
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

func GetHash(value string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(value)))
}
