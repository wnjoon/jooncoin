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

## Transaction

### UTXO

In this blockchain, We adjust UTXO(Unspent Transaction Output) model to transfer transaction.  
Transaction can be seperated into 2 models, Transaction input and Transaction output.  
- Transaction Input (TxIn) : Amount of money before user start transaction
- Transaction Output (TxOut) : Amount of money own by all users after transaction
- If User1 has $5 and give $3 to User2 who doesn't have any dollars, TxIn and TxOut might be
	- TxIn["$5(User1)"] -> We cannot divide $5. It is just one piece of paper
	- TxOut["$3(User2)", "$2(User1)", ] -> User2 should give change

So If user wants to know how much money in blockchain network, <u>Only counting on TxOut is necessary</u> and this is a mechanism of bitcoin(and also this mimic coin).  

### Coinbase

It indicates <u>how to reward to miner</u> and what to write for creation TxIn.  
It is created not by User, only by blockchain itself. And TxOut goes to miner.  

> Blockchain makes money to reward miner who verify transactions to create block in blockchain!

 


