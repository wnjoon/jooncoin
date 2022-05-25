package main

import (
	"fmt"

	"github.com/wnjoon/jooncoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockchain()
	chain.AddBlock("Second Block")
	chain.AddBlock("Thrid Block")
	chain.AddBlock("Fourth Block")

	for _, block := range chain.AllBlocks() {
		fmt.Printf("Data: %s\n", block.GetData())
		fmt.Printf("Hash: %s\n", block.GetHash())
		fmt.Printf("Prev Hash: %s\n", block.GetPrevHash())
	}
}
