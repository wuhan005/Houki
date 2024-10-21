// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"github.com/glebarez/sqlite"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var Modules ModulesStore

func New() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("houki.db"))
	if err != nil {
		return nil, errors.Wrap(err, "open database")
	}

	if err := db.AutoMigrate(
		&Module{},
	); err != nil {
		return nil, errors.Wrap(err, "auto migrate")
	}

	Modules = NewModulesStore(db)

	return db, nil
}
