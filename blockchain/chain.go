package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
			fmt.Printf("Initiaized value\n-latestHash : %s\n-Height : %d\n", b.LastHash, b.Height)
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				fmt.Println("Restoring...)")
				b.restore(checkpoint)
			}
		})
		fmt.Printf("Restoring complete\n-latestHash : %s\n-Height : %d\n", b.LastHash, b.Height)
	}
	return b
}

// Restore blockchain from database
func (b *blockchain) restore(data []byte) {
	err := gob.NewDecoder(bytes.NewReader(data)).Decode((b))
	utils.HandleError(err)
}

// Add block to blockchain
// Save blockchain updated from latest block to database for persistence
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.LastHash, b.Height+1)
	b.LastHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}
