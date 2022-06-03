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
	TxId  string `json:"txId"`  // Transaction Id used to create transaction output
	Index int    `json:"index"` // Index where transaction output is come from in transaction input
	Owner string `json:"owner"`
}

// Struct for transaction output already spent
type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

// Struct for transaction output unspent yet
type UTxOut struct {
	TxId   string
	Index  int
	Amount int
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
		// Initialize TxIn with empty value and coinbase
		{"", -1, coinbaseName},
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
	if BalanceOfAddressByTxOut(from, Blockchain()) < amount {
		return nil, errors.New("not enough amount to spend")
	}

	var txOuts []*TxOut
	var txIns []*TxIn

	// Total amount of used for spend
	total := 0

	// Get all Unspent transaction outputs of 'from address'
	uTxOuts := UTxOutsByAddress(from, Blockchain())

	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxId, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		// Sum amount of 'from address'
		total += uTxOut.Amount
	}

	// Check out change
	change := total - amount

	if change > 0 {
		// Make Transaction Out of change of 'from address'
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	// Set Transaction output of 'to address' with amount to receive
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	// Create Transaction with input/output
	tx := &Tx{
		Id:        "",
		Timestamp: utils.TimeStamp(),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}

	tx.getId()
	return tx, nil

}

// Add Transaction to mempool
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := createTx(testAddress, to, amount)
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

// Check Transaction output is still in mempool
// If it is exists, DON'T SPEND AMOUNT DOUBLE. We have to block it
func isOnMempool(uTxOut *UTxOut) bool {
	for _, tx := range Mempool.Txs {
		for _, txIn := range tx.TxIns {
			if uTxOut.TxId == txIn.TxId && uTxOut.Index == txIn.Index {
				return true
			}
		}
	}
	return false
}
