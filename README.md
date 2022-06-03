# jooncoin
Practice making coin based blockchain using golang learned by NomadCoder

> [Nomad Coder - nomad coin](https://nomadcoders.co/nomadcoin/lobby)

<br>

## Setting Environment

- go.mod : Pretty simillar to package.json in Javascript
	- go mod init github.com/wnjoon/jooncoin

## Dependencies

1. [gorilla/mux](https://github.com/gorilla/mux)
- Package gorilla/mux implements a request router and dispatcher for matching incoming requests to their respective handler.
- <u>To Use Pattern of input parameter from http request (2 line)</u>
```
r := mux.NewRouter()
r.HandleFunc("/products/{key}", ProductHandler)
r.HandleFunc("/articles/{category}/", ArticlesCategoryHandler)
r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
```
- go get -u github.com/gorilla/mux

2. [boltdb/bolt](https://github.com/boltdb/bolt)
- An embedded key/value database for Go.
- This repository has been archived by the owner. It is now read-only.
	- <u>Never be modified since it is completey developed == Stable </u>
- <i>In this program, structure would be "Hash/Block{data, hash, prevHash ....}"</i>
- In bolt, <u>Tables are called "Buckets"</u>
- go get github.com/boltdb/bolt/...

3. [evnix/boltdbweb](https://github.com/evnix/boltdbweb)
- A simple web based boltdb GUI Admin panel.
- Available to view and <u>change</u> data in database
- **Usage** : boltdbweb --db-name=<DBfilename>[required] --port=<port>[optional] --static-path=<static-path>[optional]
- go get github.com/evnix/boltdbweb

4. [br0xen/boltbrowser](https://github.com/br0xen/boltbrowser)
- A CLI Browser for BoltDB Files
- Only can view data in database
- **Usage** : boltbrowser <filename>
- go get github.com/br0xen/boltbrowser

## Proof of Work

In this blockchain, We adjust PoW(Proof of Work) consensus for approval and security used in bitcoin.  
Acutally, PoW is a little bit older strategy since its energy consumption, however, we can acknowledge inside of blockchain system such as difficulty, nonce, transactions, etc.  

- difficulty : Determine zero numbers in front of hash
- nonce : Only changable inside of blockchain

> Hard to calculate the answer(with difficulty), but easy to verify the answer.

## Used and Learned

### 1. UTXO

In this blockchain, We adjust UTXO(Unspent Transaction Output) model to transfer transaction.  
Transaction can be seperated into 2 models, Transaction input and Transaction output.  
- Transaction Input (TxIn) : Amount of money before user start transaction
- Transaction Output (TxOut) : Amount of money own by all users after transaction
- If User1 has $5 and give $3 to User2 who doesn't have any dollars, TxIn and TxOut might be
	- TxIn["$5(User1)"] -> We cannot divide $5. It is just one piece of paper
	- TxOut["$3(User2)", "$2(User1)", ] -> User2 should give change

So If user wants to know how much money in blockchain network, <u>Only counting on TxOut is necessary</u> and this is a mechanism of bitcoin(and also this mimic coin).  

### 2. Coinbase

It indicates <u>how to reward to miner</u> and what to write for creation TxIn.  
It is created not by User, only by blockchain itself. And TxOut goes to miner.  

> Blockchain makes money to reward miner who verify transactions to create block in blockchain!

### 3. Mempool

Array or Slice of unconfirmed transactions.  
Transactions in mempool don't exist in blockchain. So transactions are confirmed when miner makes a block with transaction and append block to blockchain.  
Mempool only exists in memory, not database.

> When Miner makes a block, User pay a fee to miner for their effort

### 4. Unspent Transaction Output (UTXO)

In transaction allocated in block, There are 2 types of transactions 'TxIn', 'TxOut'.  

```go
func (b *blockchain) UTxOutsByAddress(address string) []*UTxOut {
	var uTxOuts []*UTxOut

	// Find Transaction Ids used for Transaction Input from Unspent Transaction Outputs
	// User can use only from Unspent Transaction Outputs to make Transaction Inputs
	usedUtxoTxIds := make(map[string]bool)

	// 1, Search all blocks
	for _, block := range b.Blocks() {

		// 2. Search all transactions inside block (too many loop but no way)
		for _, tx := range block.Transactions {

			// 3. Check Transaction inputs to find owner is equal to address
			//    -> Transaction Input : Way point to find transaction ouputs
			for _, txIn := range tx.TxIns {
				if txIn.Owner == address {
					// If there is a TxIn from address, switch to true (it means "used!")
					usedUtxoTxIds[txIn.Owner] = true
				}
			}

			// 4. Only append UNSPENT transaction outputs to use
			for index, txOut := range tx.TxOuts {
				if txOut.Owner == address {

					// 5. Check out from usedUtxoTxIds is false 
					//    --> It means Transaction output is unspent
					_, used := usedUtxoTxIds[tx.Id]
					if !used { // if unspent(unused) => append
						uTxOuts = append(uTxOuts, &UTxOut{tx.Id, index, txOut.Amount})
					}
				}
			}
		}
	}
	return uTxOuts
}
```

### 5. Labels

If programmer wants to break a loop over than it exists, Golang supports to use 'Label'.  

```go
func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:	
	for _, tx := range Mempool.Txs {
// Outer: <- label could be written in here!
		for _, txIn := range tx.TxIns {
			if uTxOut.TxId == txIn.TxId && uTxOut.Index == txIn.Index {
				exists = true
				break Outer // Break all loops even though it is written inside a deepest loop.
			}
		}
	}
	return false
}
```

### 6. Method(Reciever function) / Function

- Method : <u>Changes</u> inside of struct value
- Function : <u>Don't change</u> inside of struct value

```go
// Method example
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := createTx(testAddress, to, amount)
	if err != nil {
		return err
	}
	// Change struct -> add transaction into mempool struct
	m.Txs = append(m.Txs, tx)	
	return nil
}
```

```go
// Function example
func Blocks(b *blockchain) []*Block {
	var blockchain []*Block
	hash := b.LastHash	// Only use for set value

	for hash != "" {
		block, _ := GetBlock(hash)
		blockchain = append(blockchain, block)
		hash = block.PrevHash
	}
	return blockchain
}
```