package kv

import (
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

type KV struct {
	db *badger.DB
}

func NewBadgerDb(pathToDb string) (*KV, error) {
	opts := badger.DefaultOptions(pathToDb)
	opts.Logger = nil // 禁用日志
	badgerInstance, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("opening kv: %w", err)
	}
	return &KV{db: badgerInstance}, nil
}

func (k *KV) Close() error {
	return k.db.Close()
}

func (k *KV) Exists(key string) (bool, error) {
	var exists bool
	err := k.db.View(func(tx *badger.Txn) error {
		_, err := tx.Get([]byte(key))
		if err == nil {
			exists = true
		}
		return err
	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		err = nil
	}
	return exists, err
}

func (k *KV) Get(key string) (string, error) {
	var value string
	err := k.db.View(func(tx *badger.Txn) error {
		item, err := tx.Get([]byte(key))
		if err != nil {
			return fmt.Errorf("getting value: %w", err)
		}
		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("copying value: %w", err)
		}
		value = string(valCopy)
		return nil
	})
	return value, err
}

func (k *KV) Set(key, value string) error {
	return k.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
}

func (k *KV) Delete(key string) error {
	return k.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}
