package blockchain

import (
	"strings"

	"github.com/wnjoon/jooncoin/db"
	"github.com/wnjoon/jooncoin/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	TimeStamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

// Create block from input parameters
// data : data to save in block
// prevHash : hash of previous block
// height : height of block
func createBlock(prevHash string, height, difficulty int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty,
		Nonce:      0,
	}
	block.mine()
	block.Transactions = Mempool.TxConfirm()
	block.persist()
	return block
}

// Get Block from input hash parameter
func GetBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, utils.ErrBlockNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

// Find certain number of zeros in front of hash for mining
func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty) // make number of zeros to verify mining
	for {
		b.TimeStamp = utils.TimeStamp()
		hash := utils.Hash(b)
		// fmt.Printf("\n\n\nTarget:%s\nHash:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			return
		}
		b.Nonce++
	}
}

/*
 * Minor Utilities
 *
 */
// Save block in database for persistence
// Shold change block struct to []byte
func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

// Restore blockhain from database
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}
