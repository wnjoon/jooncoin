package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

var host string = "http://localhost"

// Handling Error from err
func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Generate hash value from string
func GetHash(value string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(value)))
}

// Generate Bytes from interface
// interface : available to get any type of input parameter
func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)

	err := encoder.Encode(i)
	HandleError(err)

	return aBuffer.Bytes()
}

// Get data from bytes
func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(i)
	HandleError(err)
}

/*
 *
 * Minor priority utilities
 */
// Print connection information using port
func PrintConnectionInformation(port string) {
	fmt.Println("Listening on " + host + port)
}
