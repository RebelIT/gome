package gome_badger

import (
	"github.com/dgraph-io/badger"
	"log"
)

var (
	dataDirectory = "/tmp/badger"
)

type Database interface{
	Set(key,value string) error
	Delete(key string) error
	DeleteAll()error
	Get(key string) (string, error)
	GetAllKeys() ([]string, error)
	Close() error
}

type Badger struct {
	db *badger.DB
}

func Open(path string) (Database, error) {
	if path == "" {
		path = dataDirectory
	}

	d := Badger{}
	opts := badger.DefaultOptions(path)
	db, err := badger.Open(opts)
	if err != nil {
		log.Printf("[ERROR] db new %s", err)
		return d, err
	}
	d.db = db
	return d, nil
}

func (d Badger) Close() error {
	err := d.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d Badger) Set(key string, value string) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		txn.Set([]byte(key), []byte(value))
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] db set %s", err)
		return err
	}
	return nil
}

func (d Badger) Delete(key string) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		txn.Delete([]byte(key))
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] db del %s", err)
		return err
	}
	return nil
}

func (d Badger) DeleteAll() error {
	err := d.db.DropAll()
	if err != nil {
		log.Printf("[ERROR] db purge %s", err)
		return err
	}

	return nil
}

func (d Badger) Get(key string) (string, error) {
	var valCopy []byte
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			log.Printf("[ERROR] db get %s", err)
			return err
		}

		err = item.Value(func(val []byte) error {
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] db get %s", err)
			return err
		}

		valCopy, err = item.ValueCopy(nil)

		return nil
	})
	if err != nil {
		return "", err
	}
	return string(valCopy), nil
}

func (d Badger) GetAllKeys() ([]string, error) {
	keys := []string{}
	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			keys = append(keys, string(k))
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] db getAll %s", err)
		return keys, err
	}
	return keys, nil
}