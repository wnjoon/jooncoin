// package blockchain

// import (
// 	"crypto/sha256"
// 	"fmt"
// 	"sync"

// 	"github.com/wnjoon/jooncoin/utils"
// )

// type Block struct {
// 	Data     string `json:"data"`
// 	Hash     string `json:"hash"`
// 	PrevHash string `json:"previousHash,omitempty"`
// 	Height   int    `json:"height"`
// }

// type blockchain struct {
// 	blocks []*Block
// }

// // Singleton pattern
// var b *blockchain

// // To execute something just one time in entire program
// var once sync.Once

// // return type of blockchain pointer
// // for singleton pattern
// func GetBlockchain() *blockchain {
// 	if b == nil {
// 		once.Do(func() { // run only one time for initialize
// 			b = &blockchain{}
// 			b.AddBlock("Genesis block")
// 		})
// 	}
// 	return b
// }

// // data + previous hash => hash
// func (b *Block) calculateHash() {
// 	b.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(b.Data+b.PrevHash)))
// }

// func getLastHash() string {
// 	blockchainLength := len(GetBlockchain().blocks)

// 	if blockchainLength == 0 {
// 		return ""
// 	}

// 	return GetBlockchain().blocks[blockchainLength-1].Hash
// }

// func createBlock(data string) *Block {
// 	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().blocks) + 1}
// 	newBlock.calculateHash()
// 	return &newBlock
// }

// func (b *blockchain) AddBlock(data string) {
// 	b.blocks = append(b.blocks, createBlock((data)))
// }

// func (b *blockchain) AllBlocks() []*Block {
// 	return b.blocks
// }

// func (b *blockchain) GetBlockByHeight(height int) (*Block, error) {
// 	if height > len(b.blocks) {
// 		return nil, utils.ErrBlockNotFound
// 	}
// 	return b.blocks[height-1], nil
// }
