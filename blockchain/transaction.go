package blockchain

import (
	"errors"

	"github.com/wnjoon/jooncoin/utils"
)

const (
	// TODO : THIS IS A TEST VALUE => SHOULD REMOVE IT
	mineReward   int    = 50
	coinbaseName string = "COINBASE"
	testAddress  string = "joon"
)

// Struct for transaction
type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

// Struct for transaction input
type TxIn struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// Struct for transaction output
type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// Struct for mempool
type mempool struct {
	Txs []*Tx
}

// Initialize variable for mempool
var Mempool *mempool = &mempool{}

// Get Hash of Transaction Id
func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

// Create Transaction from Coinbase
// When miner creates block, coinbase gives amount of money to miner
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

// Create Transaction
// Not by coinbase (from user)
func createTx(from, to string, amount int) (*Tx, error) {
	// Check from address has enough amount
	if Blockchain().BalanceOfAddressByTxOut(from) < amount {
		return nil, errors.New("not enough amount to send")
	}

	var txIns []*TxIn
	var txOuts []*TxOut

	totalAmount := 0

	prevTxOuts := Blockchain().TxOutByAddress(from)
	for _, txOut := range prevTxOuts {
		if totalAmount > amount {
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		totalAmount += txOut.Amount
	}

	// If change exists
	change := totalAmount - amount
	if change > 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		Id:        "",
		Timestamp: utils.TimeStamp(),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	// fmt.Println(tx.Id, " ", tx.Timestamp, " ", tx.TxIns, " ", tx.TxOuts)

	return tx, nil
}

// Add Transaction to mempool
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := createTx(testAddress, to, amount)
	// fmt.Println(tx.Id, " ", tx.Timestamp, " ", tx.TxIns, " ", tx.TxOuts)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

// Confirm transaction in Mempool
// This will be happened when mining
func (m *mempool) TxConfirm() []*Tx {
	coinbase := createCoinbaseTx(testAddress)
	txs := m.Txs
	txs = append(txs, coinbase) // Add coinbase for reward of mining
	m.Txs = nil                 // delete all transactions in mempool
	return txs
}
