package provider

import (
	badger "github.com/dgraph-io/badger"
	"github.com/google/wire"
)

type BadgerDB struct {
	db      *badger.DB
	options badger.Options
}

// ProvideDatabase returns a connected badgerDB.
func ProvideDatabase(options badger.Options) (*BadgerDB, error) {
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}
	return &BadgerDB{
		db:      db,
		options: options,
	}, nil
}

// DatabaseSet is providing you with a set of database providers
var DatabaseSet = wire.NewSet(ProvideDatabase)
