package blockchain

import (
	"strconv"

	"github.com/wnjoon/jooncoin/db"
	"github.com/wnjoon/jooncoin/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash"`
	Height   int    `json:"height"`
}

// Create block from input parameters
// data : data to save in block
// prevHash : hash of previous block
// height : height of block
func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + strconv.Itoa(block.Height)
	block.Hash = utils.GetHash(payload)
	block.persist()
	return block
}

// Save block in database for persistence
// Shold change block struct to []byte
func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}
