package blockchain

import (
	"sync"

	"github.com/wnjoon/jooncoin/db"
	"github.com/wnjoon/jooncoin/utils"
)

type blockchain struct {
	LastHash string `json:"lastHash"`
	Height   int    `json:"height"`
}

var b *blockchain
var once sync.Once

// Initialize blockchain
// Load blockchain from database (blockchain.db)
// Using DB() in db/db.go
func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}

// Add block to blockchain
// Save blockchain updated from latest block to database for persistence
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.LastHash, b.Height+1)
	b.LastHash = block.Hash
	b.Height = block.Height
	b.persist()
}

// Get blockchain includes all blocks
// From database
func (b *blockchain) Blocks() []*Block {
	var blockchain []*Block
	hash := b.LastHash

	for hash != "" {
		block, _ := GetBlock(hash)
		blockchain = append(blockchain, block)
		hash = block.PrevHash
	}
	return blockchain
}

/*
 * Minor Utilities
 *
 */
// Persist block to blockchain database
func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

// Restore blockchain from database
func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}
