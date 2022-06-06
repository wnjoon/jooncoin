package blockchain

import (
	"github.com/wnjoon/jooncoin/utils"
	"github.com/wnjoon/jooncoin/wallet"
)

const (
	// TODO : THIS IS A TEST VALUE => SHOULD REMOVE IT
	mineReward   int    = 50
	CoinbaseName string = "COINBASE"
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
	TxId      string `json:"txId"`  // Transaction Id used to create transaction output
	Index     int    `json:"index"` // Index where transaction output is come from in transaction input
	Signature string `json:"signature"`
}

// Struct for transaction output already spent
type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
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
		{"", -1, CoinbaseName},
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
		return nil, utils.ErrNotEnoughMoney
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
	tx.sign()
	if !validate(tx) {
		return nil, utils.ErrNotValidated
	}
	return tx, nil

}

// Add Transaction to mempool
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := createTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

// Confirm transaction in Mempool
// This will be happened when mining
func (m *mempool) TxConfirm() []*Tx {
	coinbase := createCoinbaseTx(wallet.Wallet().Address)
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

// Sign transaction Id with wallet
func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.Id, wallet.Wallet())
	}
}

// Validate transaction
func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		// 이전에 존재하는 트랜잭션들 중에서 이번 트랜잭션에 사용된(TxIn) ID와 동일한 내용이 실제로 존재하는지 확인
		prevTx := FindTx(Blockchain(), txIn.TxId)
		if prevTx == nil {
			valid = false
			break
		}

		// If previous transaction including txId exists,
		// 이전에 존재하던 그 트랜잭션은 진짜 '내가' 만들었던 트랜잭션이 맞는지 확인
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.Id, address)
		if !valid {
			break
		}
	}

	return valid
}
