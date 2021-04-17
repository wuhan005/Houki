// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

var Modules modules

func Initialize() error {
	db, err := badger.Open(badger.DefaultOptions("./data"))
	if err != nil {
		return errors.Wrap(err, "open database")
	}

	Modules = modules{db}
	return nil
}
