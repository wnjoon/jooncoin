package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/wnjoon/jooncoin/utils"
)

// Test
const (
	fileName   string = "joon.wallet"
	walletPath string = "wallet/keys/"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

// Create private key using ecdsa
func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return privateKey
}

// Persist(save) key to file
func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)

	err = os.WriteFile(walletPath+fileName, bytes, 0644)
	utils.HandleError(err)
}

// Restore key from file
func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := os.ReadFile(walletPath + fileName)
	utils.HandleError(err)

	key, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleError(err)

	return key
}

// Get Address from privatekey
func getAddressFromPrivateKey(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

// Sign payload with privatekey
// return signature
func Sign(payload string, w *wallet) string {
	payloadBytes := decodeToBytes(payload)

	// Sign -> r, s will be created and also error
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadBytes)
	utils.HandleError(err)

	return encodeBigInts(r.Bytes(), s.Bytes())
}

// Verify signature
func verify(signature, payload, address string) bool {

	// Restore signature
	r, s, err := restoreBigInts(signature)
	utils.HandleError(err)

	// Restore address(Public key)
	x, y, err := restoreBigInts(address)
	utils.HandleError(err)

	// Get Public Key
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	payloadBytes := decodeToBytes(payload)

	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok

}

// Find a wallet file that has user's key
func hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err) // Weird function...
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			// if file exists -> restore
			w.privateKey = restoreKey()
		} else {
			privateKey := createPrivateKey()
			persistKey(privateKey)
			w.privateKey = privateKey
		}
		w.Address = getAddressFromPrivateKey(w.privateKey)
	}
	return w
}

// Get(Restore) Big Int from bytes type of input parameters
func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	payloadBytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}

	frontBytes := payloadBytes[:len(payloadBytes)/2]
	endBytes := payloadBytes[len(payloadBytes)/2:]

	bigFront, bigEnd := big.Int{}, big.Int{}
	bigFront.SetBytes(frontBytes)
	bigEnd.SetBytes(endBytes)

	return &bigFront, &bigEnd, nil
}

// Encode 2 bytes type of bigInts to string
func encodeBigInts(a, b []byte) string {
	result := append(a, b...)
	return fmt.Sprintf("%x", result)
}

// Decode String and check errors
// Just return decoded string
func decodeToBytes(payload string) []byte {
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleError(err)
	return payloadBytes
}
