package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

var host string = "http://localhost"

// Handling Error from err
func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
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

// Get Hash of any type data
func Hash(i interface{}) string {
	str := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", hash)
}

// Get Timestamp
func TimeStamp() int {
	return int(time.Now().Unix())
}

/*
 *
 * Minor priority utilities
 */
// Print connection information using port
func PrintConnectionInformation(port string) {
	fmt.Println("Listening on " + host + port)
}

// Set Rest API port
func SetPort(_port int) string {
	return fmt.Sprintf(":%d", _port)
}
