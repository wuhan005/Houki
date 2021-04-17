package db

import (
	"bytes"

	"github.com/dgraph-io/badger/v3"
)

var (
	ModulePrefix = "mod$"
)

type modules struct {
	*badger.DB
}

func (db *modules) GetEnabled() ([]string, error) {
	moduleLists := make([]string, 0)
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(ModulePrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			moduleID := append([]byte{}, it.Item().Key()...)
			moduleLists = append(moduleLists, string(bytes.TrimPrefix(moduleID, []byte(ModulePrefix))))
		}
		return nil
	})

	return moduleLists, err
}

func (db *modules) IsEnabled(id string) bool {
	return db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(ModulePrefix + id))
		return err
	}) == nil
}

func (db *modules) Enable(id string) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(ModulePrefix+id), []byte{})
	})
}

func (db *modules) Disable(id string) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(ModulePrefix + id))
	})
}
