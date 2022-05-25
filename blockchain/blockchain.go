package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []*block
}

// Singleton pattern
var b *blockchain

// To execute something just one time in entire program
var once sync.Once

// return type of blockchain pointer
// for singleton pattern
func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() { // run only one time for initialize
			b = &blockchain{}
			b.AddBlock("Genesis block")
		})
	}
	return b
}

// data + previous hash => hash
func (b *block) calculateHash() {
	b.hash = fmt.Sprintf("%x", sha256.Sum256([]byte(b.data+b.prevHash)))
}

func getLastHash() string {
	blockchainLength := len(GetBlockchain().blocks)

	if blockchainLength == 0 {
		return ""
	}

	return GetBlockchain().blocks[blockchainLength-1].hash
}

func createBlock(data string) *block {
	newBlock := block{
		data,
		"",
		getLastHash(),
	}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock((data)))
}

func (b *blockchain) AllBlocks() []*block {
	return b.blocks
}
