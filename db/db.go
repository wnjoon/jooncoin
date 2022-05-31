package db

import (
	"github.com/boltdb/bolt"
	"github.com/wnjoon/jooncoin/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blockBuckets = "blocks"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleError(err)
		db = dbPointer
		err = db.Update(func(tx *bolt.Tx) error {
			// 2 buckets created
			// 1) save data = dataBucket
			// 2) save blocks = blocksBucket
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleError(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blockBuckets))
			return err
		})
		utils.HandleError(err)
	}
	return db
}
