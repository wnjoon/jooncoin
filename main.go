package main

import (
	"github.com/wnjoon/jooncoin/blockchain"
)

func main() {
	blockchain := blockchain.GetBlockchain()
	blockchain.AddBlock("Second Block")

}
