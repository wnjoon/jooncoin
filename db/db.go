package db

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/wnjoon/jooncoin/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blockBuckets = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

// Initialize database
// From "bolt" database

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleError(err)
		db = dbPointer
		err = db.Update(func(tx *bolt.Tx) error {
			// 2 buckets created
			// 1) save data = dataBucket -> Store key(hash), value(byte[] of block)
			// 2) save blocks = blocksBucket -> Store byte[] of block
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleError(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blockBuckets))
			return err
		})
		utils.HandleError(err)
	}
	return db
}

// Get checkpoint of blockchain
// Even though program restarts, all previous data in blockchain should be restored correctly
// THIS FUNCTION SHOULD BE CALLED WHEN BLOCKCHAIN IS INITIALIZING
func Checkpoint() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

// Save "Block" with Hash and its Data
// Should make []byte type of input data
func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBuckets))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleError(err)

	fmt.Printf("Saved block\n-Hash : %s\n-Block : %b\n", hash, data)
}

// Save "Blockchain" with its Block data
// Should make []byte type of input data
func SaveBlockchain(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleError(err)

	fmt.Printf("Saved blockchain\n-Data : %b\n", data)
}

// Get Block from bucket "blockbucket"
// By Its hash []byte
func Block(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBuckets))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}
