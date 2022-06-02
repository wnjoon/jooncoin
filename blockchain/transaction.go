package blockchain

import "github.com/wnjoon/jooncoin/utils"

const (
	// TODO : THIS IS A TEST VALUE => SHOULD REMOVE IT
	mineReward   int    = 50
	coinbaseName string = "COINBASE"
)

type Tx struct {
	Id        string   `json"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func createCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{coinbaseName, mineReward},
	}
	txOuts := []*TxOut{
		{address, mineReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: utils.TimeStamp(),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}
