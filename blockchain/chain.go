package blockchain

import (
	"sync"

	"github.com/wnjoon/jooncoin/db"
	"github.com/wnjoon/jooncoin/utils"
)

type blockchain struct {
	LastHash          string `json:"lastHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDiffuculty"`
}

var b *blockchain
var once sync.Once

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

// Initialize blockchain
// Load blockchain from database (blockchain.db)
// Using DB() in db/db.go
func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock()
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}

// Add block to blockchain
// Save blockchain updated from latest block to database for persistence
func (b *blockchain) AddBlock() {
	block := createBlock(b.LastHash, b.Height+1)
	b.LastHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
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

// Get all transactionOut from blockchain
// Actually.. Due to this function, System should loop double times to get something from txOuts
// I think this function should not be exists then combine with some functions to get deterministic value
func (b *blockchain) txOuts() []*TxOut {
	var txOuts []*TxOut
	blockchain := b.Blocks()

	for _, block := range blockchain {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return txOuts
}

// Get Transaction Out of address
func (b *blockchain) TxOutByAddress(address string) []*TxOut {
	var txOutOfAddress []*TxOut
	txOuts := b.txOuts()

	for _, txOut := range txOuts {
		if txOut.Owner == address {
			txOutOfAddress = append(txOutOfAddress, txOut)
		}
	}
	return txOutOfAddress
}

// Get Balance of address by TxOut
// We have to use 2 functions with inside loops.. so wasty of memory...
func (b *blockchain) BalanceOfAddressByTxOut(address string) int {
	txOuts := b.TxOutByAddress(address)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
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

// Generate difficulty
func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// bitcoin recalculate difficulty to make 2016 blocks are made in 2 weeks
		// mimic this theory to recalculate difficulty
		return b.recalculateDifficulty()
	}
	return b.CurrentDifficulty

}

func (b *blockchain) recalculateDifficulty() int {
	blockchain := b.Blocks()
	newestBlock := blockchain[0]
	lastRecalculatedBlock := blockchain[difficultyInterval-1]

	// SInce Timestamp is second => make it minute
	actualTime := (newestBlock.TimeStamp / 60) - (lastRecalculatedBlock.TimeStamp / 60)
	expectedTime := difficultyInterval * blockInterval // Might be generated in this time

	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}
