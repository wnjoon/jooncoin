package blockchain

import (
	"strconv"

	"github.com/wnjoon/jooncoin/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash"`
	Height   int    `json:"height"`
}

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + strconv.Itoa(block.Height)
	block.Hash = utils.GetHash(payload)
	return block
}
