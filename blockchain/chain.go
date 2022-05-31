package blockchain

import "sync"

type blockchain struct {
	LastHash string `json:"lastHash"`
	Height   int    `json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.LastHash, b.Height+1)
	b.LastHash = block.Hash
	b.Height = block.Height
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			b.AddBlock("Genesis")
		})
	}
	return b
}
