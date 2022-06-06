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
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkpoint := db.Checkpoint()
		if checkpoint == nil {
			AddBlock(b)
		} else {
			b.restore(checkpoint)
		}
	})
	return b
}

/*
 *
 * Methods
 */

// Restore blockchain from database
func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

/*
 *
 * Functions
 */

// Get blockchain includes all blocks
// From database
func Blocks(b *blockchain) []*Block {
	var blockchain []*Block
	hash := b.LastHash

	for hash != "" {
		block, _ := GetBlock(hash)
		blockchain = append(blockchain, block)
		hash = block.PrevHash
	}
	return blockchain
}

// Add block to blockchain
// Save blockchain updated from latest block to database for persistence
func AddBlock(b *blockchain) {
	block := createBlock(b.LastHash, b.Height+1, getDifficulty(b))
	b.LastHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

// Persist block to blockchain database
func persistBlockchain(b *blockchain) {
	db.SaveBlockchain(utils.ToBytes(b))
}

// Recalculate difficulty using block variables
func (b *blockchain) recalculateDifficulty() int {
	blockchain := Blocks(b)
	newestBlock := blockchain[0]
	lastRecalculatedBlock := blockchain[difficultyInterval-1]

	// Since Timestamp is seconds => make it minutes
	actualTime := (newestBlock.TimeStamp / 60) - (lastRecalculatedBlock.TimeStamp / 60)
	expectedTime := difficultyInterval * blockInterval // Might be generated in this time

	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

// Generate difficulty
func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// bitcoin recalculate difficulty to make 2016 blocks are made in 2 weeks
		// mimic this theory to recalculate difficulty
		return b.recalculateDifficulty()
	}
	return b.CurrentDifficulty
}

// Get Unspent Transaction Outputs
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut

	// Find Transaction Ids used for Transaction Input from Unspent Transaction Outputs
	// User can use only from Unspent Transaction Outputs to make Transaction Inputs
	usedUtxoTxIds := make(map[string]bool)

	for _, block := range Blocks(b) {
		for _, prevTx := range block.Transactions {

			// Check Transaction inputs to find owner is equal to address
			// Is it satisfied with address is owner of transaction inputs?
			// 이제까지 송신자가 사용했던 모든 값(TxIn)들을 확인
			for _, prevTxIn := range prevTx.TxIns {

				// Should not include coinbase
				if prevTxIn.Signature == CoinbaseName {
					break
				}

				//
				if FindTx(b, prevTxIn.TxId).TxOuts[prevTxIn.Index].Address == address {
					usedUtxoTxIds[prevTxIn.TxId] = true
				}
			}

			// Only append UNSPENT transaction outputs to use
			// Check out from usedUtxoTxIds is false --> It means Transaction output is unspent
			// 송신자가 '아직 사용하지 않은 모든 TxOut들'을 확인
			for index, txOut := range prevTx.TxOuts {
				if txOut.Address == address {
					_, used := usedUtxoTxIds[prevTx.Id]
					if !used { // if unspent(unused)
						// Check this txOut is still in mempool
						uTxOut := &UTxOut{prevTx.Id, index, txOut.Amount}
						if !isOnMempool(uTxOut) { // Not in mempool => append
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

// Get Balance of address by TxOut
// We have to use 2 functions with inside loops.. so wasty of memory...
func BalanceOfAddressByTxOut(address string, b *blockchain) int {
	uTxOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, uTxOut := range uTxOuts {
		amount += uTxOut.Amount
	}
	return amount
}

// Get Transactions inside blockchain already confirmed
func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

// Find Transaction by target transaction id
// from all transactions inside blockchain(Txs)
func FindTx(b *blockchain, target string) *Tx {
	for _, tx := range Txs(b) {
		if tx.Id == target {
			return tx
		}
	}
	return nil
}
